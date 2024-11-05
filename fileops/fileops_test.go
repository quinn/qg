package fileops

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMkdirP(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test", "nested", "dir")

	if err := MkdirP(testPath + "/file.txt"); err != nil {
		t.Errorf("MkdirP() error = %v", err)
	}

	// Check if directory was created
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Error("MkdirP() failed to create directory")
	}
}

func TestWriteAndReadFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "test content"

	// Test WriteFile
	if err := WriteFile(testFile, testContent); err != nil {
		t.Errorf("WriteFile() error = %v", err)
	}

	// Test ReadFile
	content, err := ReadFile(testFile)
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
	}

	if content != testContent {
		t.Errorf("ReadFile() = %v, want %v", content, testContent)
	}
}

func TestPrint(t *testing.T) {
	// Redirect stderr to capture output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	testStr := "test message"
	Print(testStr)

	// Restore stderr
	w.Close()
	os.Stderr = oldStderr

	// Read captured output
	var buf = make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	if output != testStr {
		t.Errorf("Print() = %v, want %v", output, testStr)
	}
}

func TestGoFmt(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	// Test non-go file
	if err := GoFmt(testFile); err != nil {
		t.Errorf("GoFmt() error = %v for non-go file", err)
	}

	// Test with DRY_RUN=true
	os.Setenv("DRY_RUN", "true")
	defer os.Unsetenv("DRY_RUN")

	testGoFile := filepath.Join(tmpDir, "test.go")
	if err := GoFmt(testGoFile); err != nil {
		t.Errorf("GoFmt() error = %v with DRY_RUN=true", err)
	}
}
