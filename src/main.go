package main

import (
	"embed"
	"flag"

	. "compiler"
	. "config"
	. "config/module"
	. "render"
)

//go:generate mkdir -p build
//go:generate cp -RL ../templates build
//go:embed build/templates
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

	context := &Context{Config: config}
	if err = GenerateTemplates(embedFileSystem, context); err != nil {
		panic(err)
	}

	if err = CompileModules(config); err != nil {
		panic(err)
	}

	if err = WriteNewVersionToFile(config.Modules, config.InDir); err != nil {
		panic(err)
	}
}
