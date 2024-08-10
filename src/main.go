package main

import (
	"embed"
	"flag"

	. "proto-share/src/compiler"
	. "proto-share/src/config"
	. "proto-share/src/module"
	. "proto-share/src/template"
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

	if err = GenerateTemplates(embedFileSystem, config); err != nil {
		panic(err)
	}

	if err = CompileModules(config); err != nil {
		panic(err)
	}

	if err = UpdateMD5Hash(config.Modules, config.InDir); err != nil {
		panic(err)
	}
}
