package renderer

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	. "github.com/uben01/proto-share/internal/config"
	. "github.com/uben01/proto-share/internal/context"
	"github.com/uben01/proto-share/internal/template"
)

var templateRoot = filepath.Join("assets", "templates")

func RenderTemplates(fileSystem fs.FS, config *Config) error {
	if len(config.Languages) == 0 {
		return fmt.Errorf("no languages found in config")
	}

	if len(config.Modules) == 0 {
		return fmt.Errorf("no modules found in config")
	}

	for languageName, language := range config.Languages {
		CTX.Language = language

		fmt.Printf("Generating templates for language: %s\n", languageName)

		languageOutputPath := filepath.Join(config.OutDir, language.SubDir)

		templateLanguageRoot := filepath.Join(templateRoot, languageName.String(), "global")
		if err := walkTemplateDir(
			fileSystem,
			templateLanguageRoot,
			languageOutputPath,
		); err != nil {
			return err
		}

		templateLanguageModuleRoot := filepath.Join(templateRoot, languageName.String(), "module")
		for _, module := range config.Modules {
			if !module.Changed {
				fmt.Printf("\tModule %s has not been changed\n", module.Name)

				if !config.ForceGeneration {
					continue
				}
			}

			CTX.Module = module

			fmt.Printf("\tGenerating templates for module: %s\n", module.Name)

			moduleOutputPath := filepath.Join(
				languageOutputPath,
				template.Must(template.ProcessTemplateRecursively(language.ModuleTemplatePath, CTX)),
			)
			if err := walkTemplateDir(
				fileSystem,
				templateLanguageModuleRoot,
				moduleOutputPath,
			); err != nil {
				return err
			}
		}
	}

	return nil
}

var walkTemplateDir = func(
	fileSystem fs.FS,
	from string,
	to string,
) error {
	return fs.WalkDir(fileSystem, from, func(templateFilePath string, file os.DirEntry, err error) error {
		// It's not required to have a template for every language on both the global and module level
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}

		if err != nil {
			return err
		}

		if file.IsDir() {
			return nil
		}

		var fileContent []byte
		fileContent, err = fs.ReadFile(fileSystem, templateFilePath)
		if err != nil {
			return err
		}

		var processedTemplate string
		processedTemplate, err = template.ProcessTemplateRecursively(string(fileContent), CTX)
		if err != nil {
			return err
		}

		dir := strings.TrimPrefix(filepath.Dir(templateFilePath), from)

		return createFileFromTemplate(processedTemplate, filepath.Join(to, dir), file.Name())
	})
}

var createFileFromTemplate = func(
	processedTemplate string,
	outputFilePath string,
	outputFileName string,
) error {
	err := os.MkdirAll(outputFilePath, os.ModePerm)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(outputFilePath, outputFileName), []byte(processedTemplate), os.ModePerm)
}
