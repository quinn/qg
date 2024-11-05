package config

import (
	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
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
