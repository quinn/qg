package main

import (
	_ "embed"
	"flag"
	"log"
	"os"
	"path"

	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/hay-kot/scaffold/app/scaffold/scaffoldrc"
	"go.quinn.io/g/appdirs"
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

	resolver := pkgs.NewResolver(map[string]string{
		"gh": "https://github.com",
	}, appdirs.CacheDir(), ".")
	ppath, err := resolver.Resolve(rootDir, []string{rootDir}, &scaffoldrc.ScaffoldRC{})
	if err != nil {
		log.Fatalf("Error resolving path: %v", err)
	}

	rootDir = ppath
	configPath := path.Join(rootDir, "g.yaml")

	// Read the YAML file
	yamlData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Parse the configuration
	cfg, err := config.ParseConfig(yamlData)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	args := flag.Args()

	if len(args) == 0 {
		fileops.Print("Available generators:\n")
		for _, gen := range cfg.Generators {
			fileops.Print("* " + gen.Name)
			var args []string
			if len(gen.Use) > 0 {
				g, err := cfg.FindGenerator(gen.Use[0])
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
	generator, err := cfg.FindGenerator(gName)
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

	// Run the generator
	if err := gen.Run(gName, gConfig); err != nil {
		log.Fatal(err)
	}
}
