package shell

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunner_Run(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	runner := New(tmpDir)

	// Test with DRY_RUN=true
	os.Setenv("DRY_RUN", "true")
	if err := runner.Run("echo test"); err != nil {
		t.Errorf("Run() error = %v with DRY_RUN=true", err)
	}
	os.Unsetenv("DRY_RUN")

	// Test actual command execution
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := runner.Run("touch " + testFile); err != nil {
		t.Errorf("Run() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("Run() failed to create test file")
	}

	// Test invalid command
	if err := runner.Run("invalid-command"); err == nil {
		t.Error("Run() should return error for invalid command")
	}
}
