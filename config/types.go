package config

import (
	"fmt"

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
func (c Config) FindGenerator(gName string) (Generator, error) {
	for _, g := range c.Generators {
		if gName == g.Name {
			return g, nil
		}
	}
	return Generator{}, fmt.Errorf("generator not found: %s", gName)
}

// ParseConfig parses YAML data into a Config struct
func ParseConfig(data []byte) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML data: %w", err)
	}
	return &config, nil
}
