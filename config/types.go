package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

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
func ParseConfig(data []byte, basePath string) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML data: %w", err)
	}

	// Load included configs
	if err := config.loadIncludedConfigs(basePath); err != nil {
		return nil, fmt.Errorf("error loading included configs: %w", err)
	}

	return &config, nil
}

// loadIncludedConfigs loads and merges configs from the Include section
func (c *Config) loadIncludedConfigs(basePath string) error {
	if len(c.Include) == 0 {
		return nil
	}

	// Store original generators
	mainGenerators := c.Generators

	// Create a map to store all generators with their namespaces
	allGenerators := make([]Generator, 0)
	allGenerators = append(allGenerators, mainGenerators...)

	// Process each included config
	for namespace, includePath := range c.Include {
		// Resolve the include path relative to the base config
		fullPath := filepath.Join(basePath, includePath, "g.yaml")

		// Read the included config file
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("error reading included config %s: %w", fullPath, err)
		}

		// Parse the included config
		includedConfig, err := ParseConfig(data, filepath.Dir(fullPath))
		if err != nil {
			return fmt.Errorf("error parsing included config %s: %w", fullPath, err)
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
