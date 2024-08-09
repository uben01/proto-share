package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "proto-share/src/param"
)

func CompileModules(param *Param) error {
	for _, module := range param.Modules {
		var languageOutArgs []string

		for _, language := range param.Languages {
			err := os.MkdirAll(filepath.Join(param.OutDir, language.SubDir, language.ModulePath, module.Name, language.ProtoOutputDir), os.ModePerm)
			if err != nil {
				return err
			}

			languageProtoOutDir := filepath.Join(param.OutDir, language.SubDir, language.ModulePath, module.Name, language.ProtoOutputDir)
			languageOutArgs = append(languageOutArgs, fmt.Sprintf("--%s=%s", language.ProtocCommand, languageProtoOutDir))
		}

		protoPathForModule := filepath.Join(param.InDir, module.Path, "*.proto")
		cmdStr := fmt.Sprintf("protoc %s -I %s %s", strings.Join(languageOutArgs, " "), param.InDir, protoPathForModule)
		cmd := exec.Command("sh", "-c", cmdStr)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(output))
			return err
		}
	}

	return nil
}
