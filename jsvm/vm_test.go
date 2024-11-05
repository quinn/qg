package jsvm

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestVM_SetConfig(t *testing.T) {
	vm := New()
	config := map[string]string{
		"key": "value",
	}

	if err := vm.SetConfig(config); err != nil {
		t.Errorf("SetConfig() error = %v", err)
	}
}

func TestVM_RunConfigFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.js")

	// Create test config file
	configContent := `
		function config(input) {
			return {
				key1: input.key1,
				key2: "value2"
			};
		}
	`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	vm := New()
	config := map[string]string{
		"key1": "value1",
	}
	if err := vm.SetConfig(config); err != nil {
		t.Fatal(err)
	}

	result, err := vm.RunConfigFile(configPath, "")
	if err != nil {
		t.Errorf("RunConfigFile() error = %v", err)
	}

	expected := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("RunConfigFile() = %v, want %v", result, expected)
	}
}

func TestVM_RunTransform(t *testing.T) {
	vm := New()
	config := map[string]string{
		"key": "value",
	}
	if err := vm.SetConfig(config); err != nil {
		t.Fatal(err)
	}

	jsFunction := `function transform(input, config) {
		return input + "-" + config.key;
	}`

	if _, err := vm.vm.RunString(jsFunction); err != nil {
		t.Fatal(err)
	}

	result, err := vm.RunTransform("transform", "test", config)
	if err != nil {
		t.Errorf("RunTransform() error = %v", err)
	}

	expected := "test-value"
	if result != expected {
		t.Errorf("RunTransform() = %v, want %v", result, expected)
	}
}

func TestExportToStringMap(t *testing.T) {
	vm := New()

	tests := []struct {
		name    string
		js      string
		want    map[string]string
		wantErr bool
	}{
		{
			name: "direct string map",
			js:   `({key: "value"})`,
			want: map[string]string{
				"key": "value",
			},
		},
		{
			name:    "invalid value type",
			js:      `({key: 123})`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := vm.vm.RunString(tt.js)
			if err != nil {
				t.Fatal(err)
			}

			got, err := exportToStringMap(val)
			if (err != nil) != tt.wantErr {
				t.Errorf("exportToStringMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("exportToStringMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
