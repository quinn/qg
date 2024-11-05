package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
)

// mockResolver implements the PathResolver interface for testing
type mockResolver struct {
	rootDir string
}

func newMockResolver(rootDir string) *mockResolver {
	return &mockResolver{rootDir: rootDir}
}

// Resolve implements the PathResolver interface
func (r *mockResolver) Resolve(path string, searchPaths []string, rc pkgs.AuthProvider) (string, error) {
	// For testing, just return the path relative to rootDir
	return filepath.Join(r.rootDir, path), nil
}

func TestConfig_FindGenerator(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		gName       string
		wantErr     bool
		wantErrText string
	}{
		{
			name: "finds existing generator",
			config: Config{
				Version: "1.0",
				Generators: []Generator{
					{
						Name: "test-gen",
						Args: []string{"arg1", "arg2"},
					},
				},
			},
			gName:   "test-gen",
			wantErr: false,
		},
		{
			name: "returns error for non-existent generator",
			config: Config{
				Version:    "1.0",
				Generators: []Generator{},
			},
			gName:       "missing-gen",
			wantErr:     true,
			wantErrText: "generator not found: missing-gen",
		},
		{
			name: "finds namespaced generator",
			config: Config{
				Version: "1.0",
				Generators: []Generator{
					{
						Name: "ns:test-gen",
						Args: []string{"arg1", "arg2"},
					},
				},
			},
			gName:   "ns:test-gen",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.config.FindGenerator(tt.gName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.FindGenerator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err.Error() != tt.wantErrText {
					t.Errorf("Config.FindGenerator() error = %v, wantErrText %v", err, tt.wantErrText)
				}
				return
			}
			if got.Name != tt.gName {
				t.Errorf("Config.FindGenerator() = %v, want %v", got.Name, tt.gName)
			}
		})
	}
}

func TestParseConfig_WithIncludes(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create included config directory
	includedDir := filepath.Join(tmpDir, "included")
	if err := os.MkdirAll(includedDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create included config file
	includedConfig := []byte(`
version: "1.0"
generators:
  - name: included-gen
    args:
      - arg1
      - arg2
`)
	if err := os.WriteFile(filepath.Join(includedDir, "g.yaml"), includedConfig, 0644); err != nil {
		t.Fatal(err)
	}

	// Create main config with include
	mainConfig := []byte(`
version: "1.0"
generators:
  - name: main-gen
    args:
      - arg1
include:
  ns: included
`)

	// Create resolver
	resolver := newMockResolver(tmpDir)

	// Test parsing config with includes
	cfg, err := ParseConfig(mainConfig, tmpDir, resolver)
	if err != nil {
		t.Fatalf("ParseConfig() error = %v", err)
	}

	// Verify main generator exists
	mainGen, err := cfg.FindGenerator("main-gen")
	if err != nil {
		t.Errorf("FindGenerator() error finding main generator: %v", err)
	}
	if mainGen.Name != "main-gen" {
		t.Errorf("Main generator name = %v, want %v", mainGen.Name, "main-gen")
	}

	// Verify included generator exists with namespace
	includedGen, err := cfg.FindGenerator("ns:included-gen")
	if err != nil {
		t.Errorf("FindGenerator() error finding included generator: %v", err)
	}
	if includedGen.Name != "ns:included-gen" {
		t.Errorf("Included generator name = %v, want %v", includedGen.Name, "ns:included-gen")
	}
	if len(includedGen.Args) != 2 {
		t.Errorf("Included generator args length = %v, want %v", len(includedGen.Args), 2)
	}
}

func TestParseConfig_WithRecursiveIncludes(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create nested directory structure
	level1Dir := filepath.Join(tmpDir, "level1")
	level2Dir := filepath.Join(level1Dir, "level2")
	if err := os.MkdirAll(level2Dir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create level 2 config
	level2Config := []byte(`
version: "1.0"
generators:
  - name: level2-gen
    args:
      - arg1
`)
	if err := os.WriteFile(filepath.Join(level2Dir, "g.yaml"), level2Config, 0644); err != nil {
		t.Fatal(err)
	}

	// Create level 1 config that includes level 2
	level1Config := []byte(`
version: "1.0"
generators:
  - name: level1-gen
    args:
      - arg1
include:
  l2: level2
`)
	if err := os.WriteFile(filepath.Join(level1Dir, "g.yaml"), level1Config, 0644); err != nil {
		t.Fatal(err)
	}

	// Create main config that includes level 1
	mainConfig := []byte(`
version: "1.0"
generators:
  - name: main-gen
    args:
      - arg1
include:
  l1: level1
`)

	// Create resolver
	resolver := newMockResolver(tmpDir)

	// Test parsing config with recursive includes
	cfg, err := ParseConfig(mainConfig, tmpDir, resolver)
	if err != nil {
		t.Fatalf("ParseConfig() error = %v", err)
	}

	// Verify all generators exist with proper namespacing
	tests := []struct {
		name     string
		genName  string
		wantName string
	}{
		{"main generator", "main-gen", "main-gen"},
		{"level 1 generator", "l1:level1-gen", "l1:level1-gen"},
		{"level 2 generator", "l1:l2:level2-gen", "l1:l2:level2-gen"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, err := cfg.FindGenerator(tt.genName)
			if err != nil {
				t.Errorf("FindGenerator() error = %v", err)
				return
			}
			if gen.Name != tt.wantName {
				t.Errorf("Generator name = %v, want %v", gen.Name, tt.wantName)
			}
		})
	}
}
