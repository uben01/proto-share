package main

import (
	. "proto-share/src/compiler"
	. "proto-share/src/module"
	. "proto-share/src/param"
	. "proto-share/src/template"
)

func main() {
	params, err := ParseParams("example/proto-share.config.yml")
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
