package main

import (
	"flag"

	. "proto-share/src/compiler"
	. "proto-share/src/module"
	. "proto-share/src/param"
	. "proto-share/src/template"
)

func main() {
	configPath := flag.String("config", "", "Path to the configuration file. If not provided, read from stdin.")
	flag.Parse()
	if *configPath == "" {
		*configPath = "/dev/stdin"
	}

	params, err := ParseParams(*configPath)
	if err != nil {
		panic(err)
	}

	modules, err := GetAllModules(params.InDir)
	if err != nil {
		panic(err)
	}
	params.Modules = modules

	err = GenerateTemplates(params, modules)
	if err != nil {
		panic(err)
	}

	err = CompileModules(modules, params)
	if err != nil {
		panic(err)
	}
}
