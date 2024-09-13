package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/uben01/proto-share/internal/config"
	. "github.com/uben01/proto-share/internal/context"
	. "github.com/uben01/proto-share/internal/language"
	. "github.com/uben01/proto-share/internal/module"
	"github.com/uben01/proto-share/internal/template"
)

func CompileModules(config *Config) error {
	if len(config.Modules) == 0 {
		return fmt.Errorf("no modules defined")
	}

	numberOfLanguages := len(config.Languages)
	if numberOfLanguages == 0 {
		return fmt.Errorf("no languages defined")
	}

	anyChanged := false
	for _, module := range config.Modules {
		if !module.Changed {
			continue
		}
		anyChanged = true

		CTX.Module = module

		languageOutArgs := make([]string, numberOfLanguages)
		for _, language := range config.Languages {
			CTX.Language = language

			languageProtoOutDir, err := prepareLanguageOutput(
				config,
				language,
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

	if !anyChanged {
		fmt.Println("No output have been generated")
	}

	return nil
}

var prepareLanguageOutput = func(
	config *Config,
	language *Language,

	mkdirAll func(string, os.FileMode) error,
) (string, error) {
	pathComponents := []string{
		config.OutDir,
		language.SubDir,
		template.Must(template.ProcessTemplateRecursively(language.ModuleCompilePath, CTX)),
	}

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
