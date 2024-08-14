package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "config"
)

var (
	mkdirAll       = os.MkdirAll
	createCommand  = exec.Command
	executeCommand = func(cmd *exec.Cmd) ([]byte, error) { return cmd.CombinedOutput() }
)

func CompileModules(config *Config) error {
	if len(config.Modules) == 0 {
		return fmt.Errorf("no modules defined")
	}

	for _, module := range config.Modules {
		var languageOutArgs []string

		if len(config.Languages) == 0 {
			return fmt.Errorf("no languages defined")
		}

		for _, language := range config.Languages {
			err := mkdirAll(
				filepath.Join(config.OutDir, language.SubDir, language.ModulePath, module.Name, language.ProtoOutputDir),
				os.ModePerm,
			)
			if err != nil {
				return err
			}

			languageProtoOutDir := filepath.Join(
				config.OutDir,
				language.SubDir,
				language.ModulePath,
				module.Name,
				language.ProtoOutputDir,
			)
			languageOutArgs = append(
				languageOutArgs,
				fmt.Sprintf("--%s=%s", language.ProtocCommand, languageProtoOutDir),
			)
		}

		protoPathForModule := filepath.Join(config.InDir, module.Path, "*.proto")
		cmdStr := fmt.Sprintf(
			"protoc %s -I %s %s",
			strings.Join(languageOutArgs, " "),
			config.InDir,
			protoPathForModule,
		)
		cmd := createCommand("sh", "-c", cmdStr)
		if cmd == nil {
			return fmt.Errorf("failed to create command: %s", cmdStr)
		}

		fmt.Printf("Running command: %s\n", cmdStr)

		output, err := executeCommand(cmd)
		if err != nil {
			fmt.Println(string(output))
			return err
		}
	}

	return nil
}
