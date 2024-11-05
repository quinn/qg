package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tmpDir := os.TempDir()
	configPath := filepath.Join(tmpDir, "g.yaml")
	configContent := `
version: "1.0"
generators:
  - name: "test"
    args: ["arg1", "arg2"]
`
	os.WriteFile(configPath, []byte(configContent), 0644)
	defer os.Remove(configPath)

	config := loadConfig(tmpDir)
	if config.Version != "1.0" || len(config.Generators) != 1