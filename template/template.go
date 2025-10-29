package template

import (
	"fmt"
	"path"
	"strings"
	"text/template"

	"go.quinn.io/g/fileops"
)

// Processor handles template processing operations
type Processor struct {
	templateDir string
	outDir      string
}

// New creates a new template processor
func New(templateDir, outDir string) *Processor {
	return &Processor{
		templateDir: templateDir,
		outDir:      outDir,
	}
}

// ProcessPath processes a template path, replacing placeholders with config values
func (p *Processor) ProcessPath(templatePath string, config map[string]string) (string, error) {
	var argName string
	var brackets bool
	var targetPath string

	for _, char := range templatePath {
		switch char {
		case '[':
			if brackets {
				return "", fmt.Errorf("unterminated open bracket: %s", templatePath)
			}
			brackets = true
		case ']':
			if !brackets {
				return "", fmt.Errorf("unexpected closing bracket in path: %s", templatePath)
			}
			brackets = false
			val, ok := config[argName]
			if !ok {
				return "", fmt.Errorf("missing config value for: %s", argName)
			}
			targetPath += val
			argName = ""
		default:
			if brackets {
				argName += string(char)
			} else {
				targetPath += string(char)
			}
		}
	}

	if brackets {
		return "", fmt.Errorf("unterminated open bracket: %s", templatePath)
	}

	targetPath = path.Join(p.outDir, targetPath)
	targetPath = strings.TrimSuffix(targetPath, ".tpl")
	return targetPath, nil
}

// ProcessFile processes a template file with the given configuration
func (p *Processor) ProcessFile(sourcePath, targetPath string, config map[string]string) error {
	// Read the template file
	tmplData, err := fileops.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("error reading template file: %w", err)
	}

	// Create the target directory if it does not exist
	if err := fileops.MkdirP(targetPath); err != nil {
		return fmt.Errorf("error creating target directory: %w", err)
	}

	var result strings.Builder
	if strings.HasSuffix(sourcePath, ".tpl") {
		// Create and execute the template
		tmpl, err := template.New("file").Parse(tmplData)
		if err != nil {
			return fmt.Errorf("error parsing template file: %w", err)
		}

		// Execute the template to a string builder
		if err := tmpl.Execute(&result, config); err != nil {
			return fmt.Errorf("error executing template: %w", err)
		}
	} else {
		result.WriteString(tmplData)
	}

	// Write the result to the target file
	if err := fileops.WriteFile(targetPath, result.String()); err != nil {
		return fmt.Errorf("error writing target file: %w", err)
	}

	return nil
}
