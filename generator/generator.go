package generator

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"go.quinn.io/g/config"
	"go.quinn.io/g/fileops"
	"go.quinn.io/g/jsvm"
	"go.quinn.io/g/shell"
	tpl "go.quinn.io/g/template"
)

// Generator handles the core generation logic
type Generator struct {
	rootDir string
	Cmd     string
	Cfg     config.Generator
}

// New creates a new generator instance
func New(cfg config.Generator, cmd string, rootDir string) Generator {
	return Generator{
		rootDir: rootDir,
		Cmd:     cmd,
		Cfg:     cfg,
	}
}

// Run executes the generator with the given name and configuration
func (g *Generator) Run(generators []Generator, gConfig map[string]string, outDir string) (map[string]string, error) {
	fileops.Print("Running generator: %s\n", g.Cfg.Name)
	fileops.Print("Args: %v\n", g.Cfg.Args)
	fileops.Print("Config: %v\n", gConfig)

	if len(g.Cfg.Use) > 0 {
		for _, gName := range g.Cfg.Use {
			g, err := Find(generators, gName)
			if err != nil {
				return nil, fmt.Errorf("[USE:%s] error finding generator: %w", gName, err)
			}

			gConfigRes, err := g.Run(generators, gConfig, outDir)
			if err != nil {
				return nil, fmt.Errorf("[USE:%s] error running generator : %w", gName, err)
			}

			for k, v := range gConfigRes {
				gConfig[k] = v
			}
		}

		return gConfig, nil
	}

	// Validate required arguments
	for _, arg := range g.Cfg.Args {
		if gConfig[arg] == "" {
			return nil, fmt.Errorf("missing argument: %s", arg)
		}
	}

	templateDir := path.Join(g.rootDir, ".g", g.Cfg.Name, "tpl")
	gConfigPath := path.Join(g.rootDir, ".g", g.Cfg.Name, "config.js")

	// Process JavaScript configuration
	vm := jsvm.New()
	if err := vm.SetConfig(gConfig); err != nil {
		return nil, err
	}

	jsConfig, err := vm.RunConfigFile(gConfigPath)
	if err != nil {
		return nil, err
	}

	// Merge JavaScript config with existing config
	for k, v := range jsConfig {
		gConfig[k] = v
	}

	// Process templates
	processor := tpl.New(templateDir, outDir)
	if err := filepath.WalkDir(templateDir, func(sourcePath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		templatePath := strings.Replace(sourcePath, templateDir+"/", "", 1)
		targetPath, err := processor.ProcessPath(templatePath, gConfig)
		if err != nil {
			return err
		}

		if err := processor.ProcessFile(sourcePath, targetPath, gConfig); err != nil {
			return err
		}

		return fileops.GoFmt(targetPath)
	}); err != nil {
		return nil, fmt.Errorf("error processing templates: %w", err)
	}

	// Process transforms
	if len(g.Cfg.Transforms) > 0 {
		fileops.Print("Running transforms.\n")
		for _, transform := range g.Cfg.Transforms {
			for jsFunction, f := range transform {
				sourcePath := path.Join(outDir, f)
				sourceData, err := fileops.ReadFile(sourcePath)
				if err != nil {
					fileops.Print("Error reading file: %s\n. Will not perform transformation %s.\n original error: %v\n", sourcePath, jsFunction, err)
					continue
				}

				result, err := vm.RunTransform(jsFunction, sourceData, gConfig)
				if err != nil {
					return nil, err
				}

				if err := fileops.WriteFile(sourcePath, result); err != nil {
					return nil, err
				}

				if err := fileops.GoFmt(sourcePath); err != nil {
					return nil, err
				}
			}
		}
	}

	// Run post-generation commands
	if len(g.Cfg.Post) > 0 {
		runner := shell.New(outDir)
		for _, post := range g.Cfg.Post {
			tmpl, err := template.New("post").Parse(post)
			if err != nil {
				return nil, fmt.Errorf("error parsing post command template: %w", err)
			}

			var cmd strings.Builder
			if err := tmpl.Execute(&cmd, gConfig); err != nil {
				return nil, fmt.Errorf("error executing post command template: %w", err)
			}

			if err := runner.Run(cmd.String()); err != nil {
				return nil, fmt.Errorf("error running post command: %w", err)
			}
		}
	}

	return gConfig, nil
}
