package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dop251/goja"
	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/hay-kot/scaffold/app/scaffold/scaffoldrc"
	"go.quinn.io/g/appdirs"
	"gopkg.in/yaml.v2"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

//go:embed js/convertCase.js
var jsConvertCase string

func print(format string, a ...any) {
	if _, err := fmt.Fprintf(os.Stderr, format, a...); err != nil {
		panic(err)
	}
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
	if os.Getenv("DRY_RUN") == "true" {
		log.Println("DRY_RUN: formatting", targetPath)
		return
	}

	if strings.HasSuffix(targetPath, ".go") {
		if _, err := exec.LookPath("gopls"); err != nil {
			log.Println("gopls not found. Skipping imports and formatting.")
		} else {
			cmd := exec.Command("goimports", "-w", targetPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatalf("Error formatting imports: %v", err)
			}
			if err := exec.Command("go", "fmt", targetPath).Run(); err != nil {
				log.Fatalf("Error formatting file: %v", err)
			}
		}
	}
}

func mkdirp(targetPath string) {
	dir := filepath.Dir(targetPath)

	if os.Getenv("DRY_RUN") == "true" {
		log.Println("DRY_RUN: creating dir", dir)
		return
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("error creating target directory: %v", err)
	}
}

func must(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func execTpl(tmpl *template.Template, targetPath string, gConfig map[string]string) {
	if os.Getenv("DRY_RUN") == "true" {
		log.Println("DRY_RUN: writing to", targetPath)
		log.Println("DRY_RUN: config", gConfig)
		return
	}
	// Create the target file
	targetFile, err := os.Create(targetPath)
	if err != nil {
		log.Fatalf("error creating target file: %v", err)
	}
	defer func() { _ = targetFile.Close() }()

	// Execute the template and write to the target file
	if err := tmpl.Execute(targetFile, gConfig); err != nil {
		log.Fatalf("error executing template: %v", err)
	}
}

func write(sourcePath, data string) {
	if os.Getenv("DRY_RUN") == "true" {
		log.Println("DRY_RUN: writing to", sourcePath)
		log.Println("DRY_RUN: data", data)
		return
	}

	if err := os.WriteFile(sourcePath, []byte(data), 0644); err != nil {
		log.Fatalf("Error writing target file: %v", err)
	}
}

func errShell(outDir, cmd string) error {
	if os.Getenv("DRY_RUN") == "true" {
		log.Println("DRY_RUN: running", cmd)
		return nil
	}

	r, err := interp.New(
		interp.Interactive(true),
		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
		interp.Dir(outDir),
	)
	if err != nil {
		return err
	}

	prog, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return err
	}
	r.Reset()
	ctx := context.Background()

	print("Running shell command: %s\n", cmd)
	return r.Run(ctx, prog)
}

func shell(outDir, cmd string) {
	must(errShell(outDir, cmd))
}

type Config struct {
	Version    string      `yaml:"version"`
	Generators []Generator `yaml:"generators"`
}

// Generator represents each generator in the generators list
type Generator struct {
	Name       string              `yaml:"name"`
	Args       []string            `yaml:"args"`
	Transforms []map[string]string `yaml:"transforms"`
	Use        []string            `yaml:"use"`
	Post       []string            `yaml:"post"`
}

// Iterate over the generators and execute the one that matches the command line arguments
func findGenerator(config Config, gName string) (generator Generator) {
	found := false

	for _, g := range config.Generators {
		if gName == g.Name {
			found = true
			generator = g
			break
		}
	}
	if !found {
		log.Fatalf("Generator not found: %s", gName)
	}

	return generator
}

func main() {
	var rootDir string
	flag.StringVar(&rootDir, "path", ".", "Target directory. Contains .g dir.")
	var outDir string
	flag.StringVar(&outDir, "out", ".", "Output directory.")
	var new bool
	flag.BoolVar(&new, "new", false, "Target a new dir for generation.")

	// Custom help message
	flag.Usage = func() {
		print("Usage of %s:\n", os.Args[0])
		print("This is a sample application to demonstrate custom help text.\n\n")
		print("Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	resolver := pkgs.NewResolver(map[string]string{
		"gh": "https://github.com",
	}, appdirs.CacheDir(), ".")
	ppath, err := resolver.Resolve(rootDir, []string{rootDir}, &scaffoldrc.ScaffoldRC{})
	if err != nil {
		log.Fatalf("Error resolving path: %v", err)
	}

	rootDir = ppath
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
			var args []string
			if len(generator.Use) > 0 {
				args = findGenerator(config, generator.Use[0]).Args
			} else {
				args = generator.Args
			}

			if len(args) > 0 {
				for _, arg := range args {
					print(" [%s]", arg)
				}
			}
			print("\n")
		}
		return
	}

	args, gName := shift(args)
	generator := findGenerator(config, gName)
	gConfig := map[string]string{
		"outDir": outDir,
	}

	if len(generator.Use) > 0 {
		for ii, gName := range generator.Use {
			g := findGenerator(config, gName)

			if ii == 0 {
				// TODO: This code is duplicated.
				for i, arg := range g.Args {
					gConfig[arg] = args[i]
				}
			}

			for k, v := range runGenerator(rootDir, outDir, g, gName, gConfig) {
				gConfig[k] = v
			}
		}
	} else {
		if len(args) < len(generator.Args) {
			log.Fatalf("Missing arguments: %v", generator.Args)
		}
		// TODO: This code is duplicated.
		for i, arg := range generator.Args {
			gConfig[arg] = args[i]
		}

		runGenerator(rootDir, outDir, generator, gName, gConfig)
	}
}

func runGenerator(rootDir, outDir string, generator Generator, gName string, gConfig map[string]string) map[string]string {
	print("Running generator: %s\n", generator.Name)
	print("Args: %v\n", generator.Args)
	print("Config: %v\n", gConfig)

	for _, arg := range generator.Args {
		if gConfig[arg] == "" {
			log.Fatalf("Missing argument: %s", arg)
		}
	}

	templateDir := path.Join(rootDir, ".g", gName, "tpl")
	gConfigPath := path.Join(rootDir, ".g", gName, "config.js")

	print("Creating goja context.")
	vm := goja.New()
	print(".done\n")
	if err := vm.Set("G_CONFIG_INPUT", gConfig); err != nil {
		log.Fatalf("Error setting config input: %v", err)
	}
	print("Set config input.\n")

	// Read the config file
	print("Reading config file: %s\n", gConfigPath)
	configData, err := os.ReadFile(gConfigPath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("Error reading config file: %v", err)
		}
	} else {
		if _, err := vm.RunString(jsConvertCase); err != nil {
			log.Fatalf("Error running convertCase.js: %v", err)
		}
		if _, err := vm.RunString(string(configData)); err != nil {
			log.Fatalf("Error running config file: %v", err)
		}
		v, err := vm.RunString("config(G_CONFIG_INPUT)")
		if err != nil {
			log.Fatalf("Error running config function: %v", err)
		}

		smap, ok := v.Export().(map[string]string)
		if !ok {
			imap, ok := v.Export().(map[string]interface{})
			if !ok {
				log.Fatalf("Error exporting config: %v", v.Export())
			}
			smap = map[string]string{}
			for k, v := range imap {
				smap[k] = v.(string)
			}
		}

		for k, v := range smap {
			gConfig[k] = v
		}

		if err := vm.Set("G_CONFIG", gConfig); err != nil {
			log.Fatalf("Error setting config in VM: %v", err)
		}

		if len(generator.Transforms) > 0 {
			print("Running transforms.\n")
			for _, transform := range generator.Transforms {
				for jsFunction, f := range transform {
					sourcePath := path.Join(outDir, f)
					sourceFileData, err := os.ReadFile(sourcePath)
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
					write(sourcePath, v.String())
					gofmt(sourcePath)
				}
			}
		}
	}

	print("Config: %v\n", gConfig)

	if err := filepath.WalkDir(templateDir, func(sourcePath string, d os.DirEntry, e error) error {
		print("Walking: %s\n", sourcePath)
		if e != nil {
			return e
		}
		if d.IsDir() {
			return nil
		}
		templatePath := strings.Replace(sourcePath, templateDir+"/", "", 1)
		print("Processing template: %s\n", templatePath)

		// TODO: This could be replaced with go template by changing delims
		var argName string
		var brackets bool
		var targetPath string
		for _, char := range templatePath {
			switch char {
			case '[':
				brackets = true
			case ']':
				brackets = false
				targetPath += gConfig[argName]
				argName = ""
			default:
				if brackets {
					argName += string(char)
				} else {
					targetPath += string(char)
				}
			}
		}

		if brackets {
			log.Fatalf("unterminated open bracket: %s", templatePath)
		}

		targetPath = path.Join(outDir, targetPath)
		targetPath = strings.TrimSuffix(targetPath, ".tpl")
		print("Source path: %s\n", sourcePath)
		print("Target path: %s\n", targetPath)

		print("Creating template.\n")
		// Read the template file
		tmplData, err := os.ReadFile(sourcePath)
		if err != nil {
			return fmt.Errorf("error reading template file: %v", err)
		}

		print("Parsing template.\n")
		// Create and execute the template
		tmpl, err := template.New("file").Parse(string(tmplData))
		if err != nil {
			return fmt.Errorf("error parsing template file: %v", err)
		}

		// Create the target directory if it does not exist
		print("Making directory: %s\n", targetPath)
		mkdirp(targetPath)
		print("Executing template.\n")
		execTpl(tmpl, targetPath, gConfig)
		print("running gofmt.\n")
		gofmt(targetPath)

		return nil
	}); err != nil {
		log.Fatalf("Error walking directory: %v", err)
	}

	if len(generator.Post) > 0 {
		for _, post := range generator.Post {
			tmpl, err := template.New("file").Parse(post)
			if err != nil {
				log.Fatalf("error parsing template file: %v", err)
			}

			// exec the template against gConfig
			writeStr := strings.Builder{}
			if err := tmpl.Execute(&writeStr, gConfig); err != nil {
				log.Fatalf("error executing template: %v", err)
			}

			shell(outDir, writeStr.String())
		}
	}

	return gConfig
}
