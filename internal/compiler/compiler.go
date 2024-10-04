package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	. "github.com/uben01/proto-share/internal/config"
	. "github.com/uben01/proto-share/internal/context"
	. "github.com/uben01/proto-share/internal/language"
	. "github.com/uben01/proto-share/internal/module"
	"github.com/uben01/proto-share/internal/template"
)

func CompileModules(config *Config) {
	if len(config.Modules) == 0 {
		log.Panic("no modules defined")
	}

	numberOfLanguages := len(config.Languages)
	if numberOfLanguages == 0 {
		log.Panic("no languages defined")
	}

	anyChanged := false
	for _, module := range config.Modules {
		if !module.Changed && !config.ForceGeneration {
			continue
		}
		anyChanged = true

		CTX.Module = module

		languageOutArgs := make([]string, numberOfLanguages)
		for _, language := range config.Languages {
			CTX.Language = language

			languageProtoOutDir := prepareLanguageOutput(
				config,
				language,
				os.MkdirAll,
			)

			languageOutArgs = append(
				languageOutArgs,
				fmt.Sprintf("--%s=%s", language.ProtocCommand, languageProtoOutDir),
			)
		}

		compileProtos(
			config,
			module,
			languageOutArgs,
			exec.Command,
			func(cmd *exec.Cmd) ([]byte, error) { return cmd.CombinedOutput() },
		)
	}

	if !anyChanged {
		log.Warn("No output have been generated")
	}
}

var prepareLanguageOutput = func(
	config *Config,
	language *Language,

	mkdirAll func(string, os.FileMode) error,
) string {
	pathComponents := []string{
		config.OutDir,
		language.SubDir,
		template.ProcessTemplateRecursively(language.ModuleCompilePath, CTX),
	}

	languageProtoOutDir := filepath.Join(pathComponents...)

	if err := mkdirAll(languageProtoOutDir, os.ModePerm); err != nil {
		log.Panicf("failed to create directory: %s", err)
	}

	return languageProtoOutDir
}

func compileProtos(
	config *Config,
	module *Module,
	languageOutArgs []string,

	execute func(string, ...string) *exec.Cmd,
	combinedOutput func(cmd *exec.Cmd) ([]byte, error),
) {
	protoPathForModule := filepath.Join(config.InDir, module.Path, "*.proto")
	cmdStr := fmt.Sprintf(
		"protoc %s -I %s %s",
		strings.Trim(strings.Join(languageOutArgs, " "), " "),
		config.InDir,
		protoPathForModule,
	)

	cmd := execute("sh", "-c", cmdStr)
	if cmd == nil {
		log.Panicf("failed to create command: %s", cmdStr)
	}

	log.Infof("Running command: %s", cmdStr)

	output, err := combinedOutput(cmd)
	if err != nil {
		log.Errorf(string(output))
		log.Panicf("failed to execute command: %s", err)
	}
}
