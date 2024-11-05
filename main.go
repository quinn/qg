package main

import (
	_ "embed"
	"flag"
	"log"
	"os"

	"go.quinn.io/g/config"
	"go.quinn.io/g/fileops"
	"go.quinn.io/g/generator"
)

//go:embed js/convertCase.js
var jsConvertCase string

func shift(slice []string) ([]string, string) {
	if len(slice) == 0 {
		return slice, ""
	}
	return slice[1:], slice[0]
}

func main() {
	var rootDir string
	flag.StringVar(&rootDir, "path", ".", "Target directory. Contains .g dir.")
	var outDir string
	flag.StringVar(&outDir, "out", ".", "Output directory.")
	var new bool
	flag.BoolVar(&new, "new", false, "Target a new dir for generation.")

	// Custom help message
	flag.Usage = func() {
		fileops.Print("Usage of %s:\n", os.Args[0])
		fileops.Print("This is a sample application to demonstrate custom help text.\n\n")
		fileops.Print("Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Parse the configuration with base path and resolver for includes
	generators, err := config.LoadIncludedConfigs(rootDir, map[string]string{
		"": rootDir,
	})
	// cfg, err := config.ParseConfig(yamlData, rootDir, resolver)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	args := flag.Args()

	if len(args) == 0 {
		fileops.Print("Available generators:\n")
		for _, gen := range generators {
			fileops.Print("* " + gen.Name)
			var args []string
			if len(gen.Use) > 0 {
				g, err := config.FindGenerator(generators, gen.Use[0])
				if err != nil {
					log.Fatalf("Error finding generator: %v", err)
				}
				args = g.Args
			} else {
				args = gen.Args
			}

			if len(args) > 0 {
				for _, arg := range args {
					fileops.Print(" [%s]", arg)
				}
			}
			fileops.Print("\n")
		}
		return
	}

	args, gName := shift(args)
	gen := generator.New(rootDir, outDir, jsConvertCase)

	gConfig := map[string]string{
		"outDir": outDir,
	}

	// Find the generator and validate arguments
	generator, err := config.FindGenerator(generators, gName)
	if err != nil {
		log.Fatal(err)
	}

	if len(args) < len(generator.Args) {
		log.Fatalf("Missing arguments: %v", generator.Args)
	}

	// Set up generator config from arguments
	for i, arg := range generator.Args {
		gConfig[arg] = args[i]
	}

	// Run the generator with resolver
	if err := gen.Run(generators, gName, gConfig); err != nil {
		log.Fatal(err)
	}
}
