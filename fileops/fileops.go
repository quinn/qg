package fileops

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Print writes to stderr with formatting
func Print(format string, a ...any) {
	if _, err := fmt.Fprintf(os.Stderr, format, a...); err != nil {
		panic(err)
	}
}

// GoFmt formats a Go file using goimports and go fmt
func GoFmt(targetPath string) error {
	if os.Getenv("DRY_RUN") == "true" {
		log.Println("DRY_RUN: formatting", targetPath)
		return nil
	}

	if !strings.HasSuffix(targetPath, ".go") {
		return nil
	}

	if _, err := exec.LookPath("gopls"); err != nil {
		log.Println("gopls not found. Skipping imports and formatting.")
		return nil
	}

	cmd := exec.Command("goimports", "-w", targetPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error formatting imports: %w", err)
	}

	if err := exec.Command("go", "fmt", targetPath).Run(); err != nil {
		return fmt.Errorf("error formatting file (%s): %w", targetPath, err)
	}

	return nil
}

// MkdirP creates a directory and all necessary parent directories
func MkdirP(targetPath string) error {
	dir := filepath.Dir(targetPath)

	if os.Getenv("DRY_RUN") == "true" {
		log.Println("DRY_RUN: creating dir", dir)
		return nil
	}

	return os.MkdirAll(dir, os.ModePerm)
}

// WriteFile writes data to a file, creating the file if it does not exist
func WriteFile(sourcePath string, data string) error {
	if os.Getenv("DRY_RUN") == "true" {
		log.Println("DRY_RUN: writing to", sourcePath)
		log.Println("DRY_RUN: data", data)
		return nil
	}

	return os.WriteFile(sourcePath, []byte(data), 0644)
}

// ReadFile reads the entire file and returns it as a string
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	return string(data), nil
}
