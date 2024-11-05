package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"go.quinn.io/g/appdirs"
	"go.quinn.io/g/config"
	"go.quinn.io/g/generator"
	"gopkg.in/yaml.v2"
)

// ParseConfig parses YAML data into a Config struct and recursively loads included configs
func ParseConfig(data []byte, basePath string) (*config.Config, error) {
	var config config.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML data: %w", err)
	}

	return &config, nil
}

// loadGenerators loads and merges configs from the Include section
func LoadGenerators(basePath string, include map[string]string) ([]generator.Generator, error) {
	var allGenerators []generator.Generator

	// Process each included config
	for namespace, includePath := range include {
		resolver := pkgs.NewResolver(map[string]string{
			"gh": "https://github.com",
		}, appdirs.CacheDir(), ".")

		// Use the resolver to get the actual path of the included config
		resolvedPath, err := resolver.Resolve(includePath, []string{basePath}, nil)
		if err != nil {
			return nil, fmt.Errorf("error resolving include path %s: %w", includePath, err)
		}

		// Read the included config file
		configPath := filepath.Join(resolvedPath, "g.yaml")
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("error reading included config %s: %w", configPath, err)
		}

		// Parse the included config, passing the resolver for nested includes
		cfg, err := ParseConfig(data, resolvedPath)
		if err != nil {
			return nil, fmt.Errorf("error parsing included config %s: %w", configPath, err)
		}

		// Namespace the generators from the included config
		for _, gen := range cfg.Generators {
			cmd := gen.Name
			if namespace != "" {
				// Prefix the generator name with the namespace
				cmd = fmt.Sprintf("%s:%s", namespace, gen.Name)
			}

			if len(gen.Use) > 0 {
				if namespace != "" {
					var use []string
					for _, u := range gen.Use {
						use = append(use, fmt.Sprintf("%s:%s", namespace, u))
					}
					gen.Use = use
				}
				g, err := generator.Find(allGenerators, gen.Use[0])
				if err != nil {
					log.Fatalf("Error finding generator: %v", err)
				}
				gen.Args = g.Cfg.Args
			}

			gen := generator.New(gen, cmd, resolvedPath)
			allGenerators = append(allGenerators, gen)
		}

		generators, err := LoadGenerators(resolvedPath, cfg.Include)
		if err != nil {
			return nil, fmt.Errorf("error loading included configs: %w", err)
		}

		allGenerators = append(allGenerators, generators...)
	}

	// return nil
	return allGenerators, nil
}
