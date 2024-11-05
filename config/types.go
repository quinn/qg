package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// PathResolver is an interface for resolving paths
type PathResolver interface {
	// Resolve resolves a path to its absolute location
	Resolve(path string, searchPaths []string, rc interface{}) (string, error)
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

// FindGenerator returns the generator with the given name
// For included configs, the name should be in the format "namespace:generator"
func (c Config) FindGenerator(gName string) (Generator, error) {
	// First check in the main config's generators
	for _, g := range c.Generators {
		if gName == g.Name {
			return g, nil
		}
	}
	return Generator{}, fmt.Errorf("generator not found: %s", gName)
}

// ParseConfig parses YAML data into a Config struct and recursively loads included configs
func ParseConfig(data []byte, basePath string, resolver PathResolver) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML data: %w", err)
	}

	// Load included configs
	return &config, nil
}

// loadIncludedConfigs loads and merges configs from the Include section
func (c *Config) loadIncludedConfigs(basePath string, resolver *pkgs.Resolver) error {
	if len(c.Include) == 0 {
		return nil
	}

	// Store original generators
	mainGenerators := c.Generators

	// Create a slice to store all generators with their namespaces
	allGenerators := make([]Generator, 0)
	allGenerators = append(allGenerators, mainGenerators...)

	// Process each included config
	for namespace, includePath := range c.Include {
		// Use the resolver to get the actual path of the included config
		resolvedPath, err := resolver.Resolve(includePath, []string{basePath}, &scaffoldrc.ScaffoldRC{})
		if err != nil {
			return fmt.Errorf("error resolving include path %s: %w", includePath, err)
		}

		// Read the included config file
		configPath := filepath.Join(resolvedPath, "g.yaml")
		data, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("error reading included config %s: %w", configPath, err)
		}

		// Parse the included config, passing the resolver for nested includes
		includedConfig, err := ParseConfig(data, resolvedPath, resolver)
		if err != nil {
			return fmt.Errorf("error parsing included config %s: %w", configPath, err)
		}

		// Namespace the generators from the included config
		for _, gen := range includedConfig.Generators {
			// Prefix the generator name with the namespace
			gen.Name = fmt.Sprintf("%s:%s", namespace, gen.Name)
			allGenerators = append(allGenerators, gen)
		}
	}

	// Update the config's generators with all namespaced generators
	c.Generators = allGenerators

	return nil
}
