package main

import (
	"fmt"
	"os"
	"os/exec"
	. "proto-share/src/languages"
)

type TemplateData struct {
	JavaVersion   int
	ProjectName   string
	Modules       []*Module
	CurrentModule *Module
}

func main() {
	optInLanguages := []Language{Java()}
	inDir := "example"
	outDir := "out"

	modules, err := getAllModules(inDir)
	if err != nil {
		panic(err)
	}

	fmt.Println("Found modules:")
	for _, module := range modules {
		fmt.Printf("Name: %s, Hash: %s, Version: %d, Path: %s\n", module.Name, module.Hash, module.Version, module.Path)
	}

	templateData := TemplateData{
		JavaVersion: 21,
		ProjectName: "myProject",
		Modules:     modules,
	}

	err = generateTemplates(optInLanguages, modules, outDir, templateData)
	if err != nil {
		panic(err)
	}

	for _, module := range modules {
		languageOutArgs := ""

		for _, language := range optInLanguages {
			err := os.MkdirAll(fmt.Sprintf("%s/%s/%s/%s/%s", outDir, language.SubDir(), language.ModulePath(), module.Name, language.ProtoOutputDir()), os.ModePerm)
			if err != nil {
				panic(err)
			}

			languageOutArgs += fmt.Sprintf(" --%s=%s/%s/%s/%s/%s", language.ProtocCommand(), outDir, language.SubDir(), language.ModulePath(), module.Name, language.ProtoOutputDir())
		}

		cmdStr := fmt.Sprintf("protoc %s -I %s %s/%s/*.proto", languageOutArgs, inDir, inDir, module.Path)
		cmd := exec.Command("sh", "-c", cmdStr)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(output))
			panic(err)
		}
	}
}
