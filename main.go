package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
	"gopkg.in/yaml.v2"
)

func print(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func shift(slice []string) ([]string, string) {
	if len(slice) == 0 {
		return slice, ""
	}
	return slice[1:], slice[0]
}

func ext(filename string) (string, string) {
	parts := strings.SplitN(filename, ".", 2)
	if len(parts) == 1 {
		return filename, ""
	}

	return parts[0], parts[1]
}

func gofmt(targetPath string) {
	if strings.HasSuffix(targetPath, ".go") {
		if _, err := exec.LookPath("gopls"); err != nil {
			log.Println("gopls not found. Skipping imports and formatting.")
		} else {
			if err := exec.Command("gopls", "imports", "-w", targetPath).Run(); err != nil {
				log.Fatalf("Error formatting imports: %v", err)
			}
			if err := exec.Command("gopls", "format", "-w", targetPath).Run(); err != nil {
				log.Fatalf("Error formatting file: %v", err)
			}
		}
	}
}

type Config struct {
	Version    string      `yaml:"version"`
	Generators []Generator `yaml:"generators"`
}

// Generator represents each generator in the generators list
type Generator struct {
	Name      string              `yaml:"name"`
	Args      []string            `yaml:"args"`
	Transorms []map[string]string `yaml:"transforms"`
}

func main() {
	var rootDir string

	flag.StringVar(&rootDir, "path", ".", "Target directory. Contains .g dir.")

	// Custom help message
	flag.Usage = func() {
		print("Usage of %s:\n", os.Args[0])
		print("This is a sample application to demonstrate custom help text.\n\n")
		print("Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	configPath := path.Join(rootDir, "g.yaml")

	// Read the YAML file
	yamlData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Unmarshal the YAML data into a Config struct
	var config Config
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		log.Fatalf("Error unmarshalling YAML data: %v", err)
	}

	args := flag.Args()

	if len(args) == 0 {
		print("Available generators:\n")
		for _, generator := range config.Generators {
			print("* " + generator.Name)
			if len(generator.Args) > 0 {
				for _, arg := range generator.Args {
					print(" [%s]", arg)
				}
			}
			print("\n")
		}
		return
	}

	args, gName := shift(args)
	found := false
	var generator Generator

	// Iterate over the generators and execute the one that matches the command line arguments
	for _, g := range config.Generators {
		if gName == g.Name {
			found = true
			generator = g
			print("Running generator: %s\n", g.Name)
			print("Args: %v\n", g.Args)
			if len(args) != len(g.Args) {
				log.Fatalf("Invalid number of arguments. Expected %v, got %d", g.Args, len(args))
			}
			break
		}
	}

	if !found {
		log.Fatalf("Generator not found: %s", gName)
	}

	gConfig := map[string]string{}

	for i, arg := range generator.Args {
		gConfig[arg] = args[i]
	}

	print("Config: %v\n", gConfig)

	templateDir := path.Join(rootDir, ".g", gName, "tpl")
	gConfigPath := path.Join(rootDir, ".g", gName, "config.js")

	vm := goja.New()
	if err := vm.Set("G_CONFIG_INPUT", gConfig); err != nil {
		log.Fatalf("Error setting config input: %v", err)
	}
	// Read the config file
	configData, err := os.ReadFile(gConfigPath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("Error reading config file: %v", err)
		}
	} else {
		if _, err := vm.RunString(string(configData)); err != nil {
			log.Fatalf("Error running config file: %v", err)
		}
		v, err := vm.RunString("config(G_CONFIG_INPUT)")
		if err != nil {
			log.Fatalf("Error running config function: %v", err)
		}
		for k, v := range v.Export().(map[string]interface{}) {
			gConfig[k] = v.(string)
		}

		if err := vm.Set("G_CONFIG", gConfig); err != nil {
			log.Fatalf("Error setting config: %v", err)
		}

		if len(generator.Transorms) > 0 {
			for _, transform := range generator.Transorms {
				for jsFunction, sourceFile := range transform {
					vm := goja.New()
					if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
						log.Fatalf("Target file does not exist: %s", sourceFile)
					}
					sourceFileData, err := os.ReadFile(sourceFile)
					if err != nil {
						log.Fatalf("Error reading target file: %v", err)
					}
					if err := vm.Set("G_FILE_INPUT", string(sourceFileData)); err != nil {
						log.Fatalf("Error setting file input: %v", err)
					}
					if _, err := vm.RunString(string(configData)); err != nil {
						log.Fatalf("Error running config file: %v", err)
					}
					v, err := vm.RunString(jsFunction + "(G_FILE_INPUT, G_CONFIG)")
					if err != nil {
						log.Fatalf("Error running transform function: %v", err)
					}
					if err := os.WriteFile(sourceFile, []byte(v.String()), 0644); err != nil {
						log.Fatalf("Error writing target file: %v", err)
					}
					gofmt(sourceFile)
				}
			}
		}
	}

	print("Config: %v\n", gConfig)

	if err := filepath.WalkDir(templateDir, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() {
			return nil
		}
		sourcePath := p
		p = strings.Replace(p, templateDir+"/", "", 1)
		print("Processing template: %s\n", p)
		targetPathComponents := []string{rootDir}
		strings.Split(p, "/")
		for _, c := range strings.Split(p, "/") {
			if strings.HasPrefix(c, "[") && strings.Contains(c, "]") {
				fn, ext := ext(c)
				c = strings.TrimSuffix(strings.TrimPrefix(fn, "["), "]")
				c = gConfig[c]
				c = c + "." + ext
			}
			targetPathComponents = append(targetPathComponents, c)
		}

		targetPath := path.Join(targetPathComponents...)
		targetPath = strings.TrimSuffix(targetPath, ".tpl")
		print("Source path: %s\n", sourcePath)
		print("Target path: %s\n", targetPath)

		// Read the template file
		tmplData, err := os.ReadFile(sourcePath)
		if err != nil {
			return fmt.Errorf("error reading template file: %v", err)
		}

		// Create and execute the template
		tmpl, err := template.New("file").Parse(string(tmplData))
		if err != nil {
			return fmt.Errorf("error parsing template file: %v", err)
		}

		// Create the target directory if it does not exist
		if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
			return fmt.Errorf("error creating target directory: %v", err)
		}

		// Create the target file
		targetFile, err := os.Create(targetPath)
		if err != nil {
			return fmt.Errorf("error creating target file: %v", err)
		}
		defer targetFile.Close()

		// Execute the template and write to the target file
		if err := tmpl.Execute(targetFile, gConfig); err != nil {
			return fmt.Errorf("error executing template: %v", err)
		}

		gofmt(targetPath)

		return nil
	}); err != nil {
		log.Fatalf("Error walking directory: %v", err)
	}
}
