package main

import (
	_ "embed"
	"flag"
	"log"
	"os"

	"go.quinn.io/g/fileops"
	"go.quinn.io/g/generator"
	"go.quinn.io/g/util"
)

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
	generators, err := util.LoadGenerators(rootDir, map[string]string{
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
			fileops.Print("* " + gen.Cmd)
			var args []string
			if len(gen.Cfg.Use) > 0 {
				g, err := generator.Find(generators, gen.Cfg.Use[0])
				if err != nil {
					log.Fatalf("Error finding generator: %v", err)
				}
				args = g.Cfg.Args
			} else {
				args = gen.Cfg.Args
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
	// gen := generator.New(rootDir, outDir, jsConvertCase)

	gConfig := map[string]string{
		"outDir": outDir,
	}

	// Find the configG and validate arguments
	gen, err := generator.Find(generators, gName)
	if err != nil {
		log.Fatal(err)
	}

	if len(args) < len(gen.Cfg.Args) {
		log.Fatalf("Missing arguments: %v", gen.Cfg.Args)
	}

	// Set up generator config from arguments
	for i, arg := range gen.Cfg.Args {
		gConfig[arg] = args[i]
	}

	if err := gen.Run(gConfig, outDir); err != nil {
		log.Fatal(err)
	}
}
