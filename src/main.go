package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	. "proto-share/src/languages"
	"strings"
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
		var languageOutArgs []string

		for _, language := range optInLanguages {
			err := os.MkdirAll(filepath.Join(outDir, language.SubDir(), language.ModulePath(), module.Name, language.ProtoOutputDir()), os.ModePerm)
			if err != nil {
				panic(err)
			}

			languageProtoOutDir := filepath.Join(outDir, language.SubDir(), language.ModulePath(), module.Name, language.ProtoOutputDir())
			languageOutArgs = append(languageOutArgs, fmt.Sprintf("--%s=%s", language.ProtocCommand(), languageProtoOutDir))
		}

		protoPathForModule := filepath.Join(inDir, module.Path, "*.proto")
		cmdStr := fmt.Sprintf("protoc %s -I %s %s", strings.Join(languageOutArgs, " "), inDir, protoPathForModule)
		cmd := exec.Command("sh", "-c", cmdStr)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(output))
			panic(err)
		}
	}
}
