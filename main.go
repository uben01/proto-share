package main

import (
	"embed"
	"flag"

	. "github.com/uben01/proto-share/internal/compiler"
	. "github.com/uben01/proto-share/internal/config"
	. "github.com/uben01/proto-share/internal/module"
	. "github.com/uben01/proto-share/internal/render"
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

	modules, err := GetAllModules(config.InDir)
	if err != nil {
		panic(err)
	}
	config.Modules = modules

	if err = UpdateModulesVersion(config.Modules, config.InDir); err != nil {
		panic(err)
	}

	if err = GenerateTemplates(embedFileSystem, config); err != nil {
		panic(err)
	}

	if err = CompileModules(config); err != nil {
		panic(err)
	}

	if err = WriteNewVersionToFile(config.Modules, config.InDir); err != nil {
		panic(err)
	}
}
