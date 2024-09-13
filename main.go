package main

import (
	"embed"
	"flag"
	"fmt"

	. "github.com/uben01/proto-share/internal/compiler"
	. "github.com/uben01/proto-share/internal/config"
	. "github.com/uben01/proto-share/internal/context"
	. "github.com/uben01/proto-share/internal/module"
	. "github.com/uben01/proto-share/internal/renderer"
)

//go:embed assets/templates
var embedFileSystem embed.FS

func main() {
	configPath := flag.String("config", "", "Path to the configuration file. If not provided, read from stdin.")
	flag.Parse()
	if *configPath == "" {
		*configPath = "/dev/stdin"
	}

	config, err := ParseConfig(*configPath)
	if err != nil {
		panic(err)
	}
	CTX.Config = config

	modules, err := DiscoverModules(config.InDir)
	if err != nil {
		panic(err)
	}
	config.Modules = modules

	if err = UpdateModuleVersions(config.Modules, config.InDir); err != nil {
		panic(err)
	}

	if !config.AnyModuleChanged() && !config.ForceGeneration {
		fmt.Println("No changes detected. Exiting.")

		return
	}

	if err = RenderTemplates(embedFileSystem, config); err != nil {
		panic(err)
	}

	if err = CompileModules(config); err != nil {
		panic(err)
	}

	if err = WriteNewVersionToFile(config.Modules, config.InDir); err != nil {
		panic(err)
	}
}
