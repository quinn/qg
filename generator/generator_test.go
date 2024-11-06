package generator

import (
	"os"
	"path/filepath"
	"testing"

	"go.quinn.io/g/config"
)

func TestGenerator_Run(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	rootDir := filepath.Join(tmpDir, "root")
	outDir := filepath.Join(tmpDir, "out")

	// Create test directory structure
	must(t, os.MkdirAll(filepath.Join(rootDir, ".g", "test-gen", "tpl"), 0755))

	// Create test template
	tplContent := "Hello {{.name}}!"
	must(t, os.WriteFile(filepath.Join(rootDir, ".g", "test-gen", "tpl", "test.txt.tpl"), []byte(tplContent), 0644))

	// Create generator instance
	cfg := config.Generator{
		Name: "test-gen",
		Args: []string{"name"},
	}
	g := New(cfg, "test-gen", rootDir)

	// Run generator
	generators := []Generator{g}
	_, err := g.Run(generators, map[string]string{
		"name": "World",
	}, outDir)
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
	cfg := config.Generator{
		Name: "test-gen",
		Args: []string{"name"},
		Transforms: []map[string]string{
			{"transform": "test.txt"},
		},
	}
	g := New(cfg, "test-gen", rootDir)

	// Run generator
	generators := []Generator{g}
	_, err := g.Run(generators, map[string]string{
		"name": "World",
	}, outDir)
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

	// Create test template
	tplContent := "Hello {{.name}}!"
	must(t, os.WriteFile(filepath.Join(rootDir, ".g", "test-gen", "tpl", "test.txt.tpl"), []byte(tplContent), 0644))

	// Create generator instance
	cfg := config.Generator{
		Name: "test-gen",
		Args: []string{"name"},
		Post: []string{"touch {{.name}}.flag"},
	}
	g := New(cfg, "test-gen", rootDir)

	// Run generator
	generators := []Generator{g}
	_, err := g.Run(generators, map[string]string{
		"name": "test",
	}, outDir)
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
