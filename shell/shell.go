package shell

import (
	"context"
	"fmt"
	"os"
	"strings"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// Runner handles shell command execution
type Runner struct {
	workDir string
}

// New creates a new shell command runner
func New(workDir string) *Runner {
	return &Runner{
		workDir: workDir,
	}
}

// Run executes a shell command in the specified working directory
func (r *Runner) Run(cmd string) error {
	if os.Getenv("DRY_RUN") == "true" {
		fmt.Printf("DRY_RUN: running %s\n", cmd)
		return nil
	}

	runner, err := interp.New(
		interp.Interactive(true),
		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
		interp.Dir(r.workDir),
	)
	if err != nil {
		return fmt.Errorf("error creating shell interpreter: %w", err)
	}

	prog, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return fmt.Errorf("error parsing shell command: %w", err)
	}

	runner.Reset()
	ctx := context.Background()

	fmt.Printf("Running shell command: %s\n", cmd)
	return runner.Run(ctx, prog)
}
