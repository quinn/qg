package jsvm

import (
	"fmt"
	"os"

	"github.com/dop251/goja"
)

// VM wraps the JavaScript virtual machine functionality
type VM struct {
	vm *goja.Runtime
}

// New creates a new JavaScript VM instance
func New() *VM {
	return &VM{
		vm: goja.New(),
	}
}

// SetConfig sets the configuration in the VM environment
func (v *VM) SetConfig(config map[string]string) error {
	if err := v.vm.Set("G_CONFIG_INPUT", config); err != nil {
		return fmt.Errorf("error setting config input: %w", err)
	}
	return nil
}

// RunConfigFile executes a JavaScript config file and returns the resulting configuration
func (v *VM) RunConfigFile(configPath string, convertCaseJS string) (map[string]string, error) {
	// Run the convertCase.js helper
	if _, err := v.vm.RunString(convertCaseJS); err != nil {
		return nil, fmt.Errorf("error running convertCase.js: %w", err)
	}

	// Read and run the config file
	configData, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if _, err := v.vm.RunString(string(configData)); err != nil {
		return nil, fmt.Errorf("error running config file: %w", err)
	}

	// Execute the config function
	result, err := v.vm.RunString("config(G_CONFIG_INPUT)")
	if err != nil {
		return nil, fmt.Errorf("error running config function: %w", err)
	}

	return exportToStringMap(result)
}

// RunTransform executes a JavaScript transform function on the given input
func (v *VM) RunTransform(jsFunction string, fileInput string, config map[string]string) (string, error) {
	if err := v.vm.Set("G_FILE_INPUT", fileInput); err != nil {
		return "", fmt.Errorf("error setting file input: %w", err)
	}

	if err := v.vm.Set("G_CONFIG", config); err != nil {
		return "", fmt.Errorf("error setting config: %w", err)
	}

	result, err := v.vm.RunString(jsFunction + "(G_FILE_INPUT, G_CONFIG)")
	if err != nil {
		return "", fmt.Errorf("error running transform function: %w", err)
	}

	return result.String(), nil
}

// exportToStringMap converts a goja.Value to a map[string]string
func exportToStringMap(v goja.Value) (map[string]string, error) {
	export := v.Export()

	// Try direct conversion to map[string]string
	if smap, ok := export.(map[string]string); ok {
		return smap, nil
	}

	// Try conversion from map[string]interface{}
	if imap, ok := export.(map[string]interface{}); ok {
		smap := make(map[string]string)
		for k, v := range imap {
			str, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("value for key %s is not a string", k)
			}
			smap[k] = str
		}
		return smap, nil
	}

	return nil, fmt.Errorf("unable to convert value to map[string]string")
}
