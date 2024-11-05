package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"go.quinn.io/g/appdirs"
	"gopkg.in/yaml.v2"
)

// PathResolver is an interface for resolving paths
type PathResolver interface {
	// Resolve resolves a path to its absolute location
	Resolve(path string, searchPaths []string, rc pkgs.AuthProvider) (string, error)
}

// Config represents the main configuration structure
type Config struct {
	Version    string            `yaml:"version"`
	Generators []Generator       `yaml:"generators"`
	Include    map[string]string `yaml:"include"`
}

// Generator represents each generator in the generators list
type Generator struct {
	Name       string              `yaml:"name"`
	Args       []string            `yaml:"args"`
	Transforms []map[string]string `yaml:"transforms"`
	Use        []string            `yaml:"use"`
	Post       []string            `yaml:"post"`
}

func FindGenerator(generators []Generator, name string) (*Generator, error) {
	for _, gen := range generators {
		if gen.Name == name {
			return &gen, nil
		}
	}
	return nil, fmt.Errorf("generator not found: %s", name)
}

// ParseConfig parses YAML data into a Config struct and recursively loads included configs
func ParseConfig(data []byte, basePath string) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML data: %w", err)
	}

	// Load included configs
	if len(config.Include) != 0 {

		// ppath, err := resolver.Resolve(rootDir, []string{rootDir}, &scaffoldrc.ScaffoldRC{})
		// if err != nil {
		// 	log.Fatalf("Error resolving path: %v", err)
		// }

		// Store original generators
		// mainGenerators := config.Generators

		// Create a slice to store all generators with their namespaces
		// allGenerators := make([]Generator, 0)
		// allGenerators = append(allGenerators, mainGenerators...)

		generators, err := LoadIncludedConfigs(basePath, config.Include)
		if err != nil {
			return nil, fmt.Errorf("error loading included configs: %w", err)
		}

		config.Generators = append(config.Generators, generators...)
	}

	return &config, nil
}

// loadIncludedConfigs loads and merges configs from the Include section
func LoadIncludedConfigs(basePath string, include map[string]string) ([]Generator, error) {
	var allGenerators []Generator
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
		includedConfig, err := ParseConfig(data, resolvedPath)
		if err != nil {
			return nil, fmt.Errorf("error parsing included config %s: %w", configPath, err)
		}

		// Namespace the generators from the included config
		for _, gen := range includedConfig.Generators {
			// Prefix the generator name with the namespace
			if namespace != "" {
				gen.Name = fmt.Sprintf("%s:%s", namespace, gen.Name)
			}
			allGenerators = append(allGenerators, gen)
		}
	}

	// // Update the config's generators with all namespaced generators
	// c.Generators = allGenerators

	// return nil
	return allGenerators, nil
}
