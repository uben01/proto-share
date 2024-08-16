package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "config"
	. "config/language"
	. "config/module"
)

func CompileModules(config *Config) error {
	if len(config.Modules) == 0 {
		return fmt.Errorf("no modules defined")
	}

	numberOfLanguages := len(config.Languages)
	if numberOfLanguages == 0 {
		return fmt.Errorf("no languages defined")
	}

	for _, module := range config.Modules {
		languageOutArgs := make([]string, numberOfLanguages)
		for _, language := range config.Languages {
			languageProtoOutDir, err := prepareLanguageOutput(
				config,
				language,
				module,
				os.MkdirAll,
			)
			if err != nil {
				return err
			}

			languageOutArgs = append(
				languageOutArgs,
				fmt.Sprintf("--%s=%s", language.ProtocCommand, languageProtoOutDir),
			)
		}

		if err := compileProtos(
			config,
			module,
			languageOutArgs,
			exec.Command,
			func(cmd *exec.Cmd) ([]byte, error) { return cmd.CombinedOutput() },
		); err != nil {
			return err
		}
	}

	return nil
}

func prepareLanguageOutput(
	config *Config,
	language *Language,
	module *Module,

	mkdirAll func(string, os.FileMode) error,
) (string, error) {
	pathComponents := []string{
		config.OutDir,
		language.SubDir,
		language.ModulePath,
	}

	if language.SeparateModuleDir {
		pathComponents = append(pathComponents, module.Name)
	}

	pathComponents = append(pathComponents, language.ProtoOutputDir)

	languageProtoOutDir := filepath.Join(pathComponents...)

	if err := mkdirAll(languageProtoOutDir, os.ModePerm); err != nil {
		return "", err
	}
	return languageProtoOutDir, nil
}

func compileProtos(
	config *Config,
	module *Module,
	languageOutArgs []string,

	execute func(string, ...string) *exec.Cmd,
	combinedOutput func(cmd *exec.Cmd) ([]byte, error),
) error {
	protoPathForModule := filepath.Join(config.InDir, module.Path, "*.proto")
	cmdStr := fmt.Sprintf(
		"protoc %s -I %s %s",
		strings.Join(languageOutArgs, " "),
		config.InDir,
		protoPathForModule,
	)

	cmd := execute("sh", "-c", cmdStr)
	if cmd == nil {
		return fmt.Errorf("failed to create command: %s", cmdStr)
	}

	fmt.Printf("Running command: %s\n", cmdStr)

	output, err := combinedOutput(cmd)
	if err != nil {
		fmt.Println(string(output))
		return err
	}

	return nil
}
