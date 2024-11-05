package generator

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

func TestGenerator_Run(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	rootDir := filepath.Join(tmpDir, "root")
	outDir := filepath.Join(tmpDir, "out")

	// Create test directory structure

	// Create test config
	configYaml := `
version: "1.0"
generators:
  - name: test-gen
    args:
      - name
    transforms: []
    post: []
`
	must(t, os.WriteFile(filepath.Join(rootDir, "g.yaml"), []byte(configYaml), 0644))

	// Create test template
	tplContent := "Hello {{.name}}!"
	must(t, os.WriteFile(filepath.Join(rootDir, ".g", "test-gen", "tpl", "test.txt.tpl"), []byte(tplContent), 0644))

	// Create generator instance
	g := New(rootDir, outDir, "")
	resolver := newMockResolver(rootDir)

	// Run generator
	err := g.Run("test-gen", map[string]string{
		"name": "World",
	}, resolver)
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}

	// Verify output
	content, err := os.ReadFile(filepath.Join(outDir, "test.txt"))
	if err != nil {
		t.Fatal(err)
	}

	expected := "Hello World!"
	if string(content) != expected {
		t.Errorf("Run() output = %v, want %v", string(content), expected)
	}
}

func TestGenerator_RunWithTransforms(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	rootDir := filepath.Join(tmpDir, "root")
	outDir := filepath.Join(tmpDir, "out")

	// Create test directory structure
	must(t, os.MkdirAll(filepath.Join(rootDir, ".g", "test-gen", "tpl"), 0755))

	// Create test config with transform
	configYaml := `
version: "1.0"
generators:
  - name: test-gen
    args:
      - name
    transforms:
      - transform: test.txt
    post: []
`
	must(t, os.WriteFile(filepath.Join(rootDir, "g.yaml"), []byte(configYaml), 0644))

	// Create test template
	tplContent := "Hello {{.name}}!"
	must(t, os.WriteFile(filepath.Join(rootDir, ".g", "test-gen", "tpl", "test.txt.tpl"), []byte(tplContent), 0644))

	// Create test config.js
	configJS := `
function config(input) {
    return {
        name: input.name.toUpperCase()
    };
}

function transform(input, config) {
    return input + " Transformed!";
}
`
	must(t, os.WriteFile(filepath.Join(rootDir, ".g", "test-gen", "config.js"), []byte(configJS), 0644))

	// Create generator instance
	g := New(rootDir, outDir, "")
	resolver := newMockResolver(rootDir)

	// Run generator
	err := g.Run("test-gen", map[string]string{
		"name": "World",
	}, resolver)
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}

	// Verify output
	content, err := os.ReadFile(filepath.Join(outDir, "test.txt"))
	if err != nil {
		t.Fatal(err)
	}

	expected := "Hello WORLD! Transformed!"
	if string(content) != expected {
		t.Errorf("Run() output = %v, want %v", string(content), expected)
	}
}

func TestGenerator_RunWithPost(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	rootDir := filepath.Join(tmpDir, "root")
	outDir := filepath.Join(tmpDir, "out")

	// Create test directory structure
	must(t, os.MkdirAll(filepath.Join(rootDir, ".g", "test-gen", "tpl"), 0755))

	// Create test config with post command
	configYaml := `
version: "1.0"
generators:
  - name: test-gen
    args:
      - name
    transforms: []
    post:
      - "touch {{.name}}.flag"
`
	must(t, os.WriteFile(filepath.Join(rootDir, "g.yaml"), []byte(configYaml), 0644))

	// Create test template
	tplContent := "Hello {{.name}}!"
	must(t, os.WriteFile(filepath.Join(rootDir, ".g", "test-gen", "tpl", "test.txt.tpl"), []byte(tplContent), 0644))

	// Create generator instance
	g := New(rootDir, outDir, "")
	resolver := newMockResolver(rootDir)

	// Run generator
	err := g.Run("test-gen", map[string]string{
		"name": "test",
	}, resolver)
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}

	// Verify post command execution
	if _, err := os.Stat(filepath.Join(outDir, "test.flag")); os.IsNotExist(err) {
		t.Error("Post command did not create flag file")
	}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
