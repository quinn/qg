package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessor_ProcessPath(t *testing.T) {
	processor := New("/templates", "/output")
	config := map[string]string{
		"name": "test",
		"type": "component",
	}

	tests := []struct {
		name        string
		path        string
		want        string
		wantErr     bool
		wantErrText string
	}{
		{
			name: "simple path",
			path: "file.txt",
			want: "/output/file.txt",
		},
		{
			name: "path with placeholders",
			path: "[type]/[name].txt",
			want: "/output/component/test.txt",
		},
		{
			name:        "unterminated bracket",
			path:        "[type/[name].txt",
			wantErr:     true,
			wantErrText: "unterminated open bracket: [type/[name].txt",
		},
		{
			name:        "missing config value",
			path:        "[missing].txt",
			wantErr:     true,
			wantErrText: "missing config value for: missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processor.ProcessPath(tt.path, config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err.Error() != tt.wantErrText {
					t.Errorf("ProcessPath() error = %v, wantErrText %v", err, tt.wantErrText)
				}
				return
			}
			if got != tt.want {
				t.Errorf("ProcessPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessor_ProcessFile(t *testing.T) {
	// Create temporary directories for testing
	tmpDir := t.TempDir()
	templateDir := filepath.Join(tmpDir, "templates")
	outDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test template file
	templateContent := "Hello {{.name}}!"
	templatePath := filepath.Join(templateDir, "test.txt.tpl")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatal(err)
	}

	processor := New(templateDir, outDir)
	config := map[string]string{
		"name": "World",
	}

	targetPath := filepath.Join(outDir, "test.txt")
	if err := processor.ProcessFile(templatePath, targetPath, config); err != nil {
		t.Errorf("ProcessFile() error = %v", err)
	}

	// Verify output
	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}

	expected := "Hello World!"
	if string(content) != expected {
		t.Errorf("ProcessFile() output = %v, want %v", string(content), expected)
	}
}
