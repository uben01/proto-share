package renderer

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	templ "text/template"

	. "github.com/uben01/proto-share/internal/config"
)

var templateRoot = filepath.Join("assets", "templates")

func RenderTemplates(fileSystem fs.FS, config *Config) error {
	context := &context{Config: config}

	for languageName, language := range config.Languages {
		context.Language = language

		fmt.Printf("Generating templates for language: %s\n", languageName)

		languageOutputPath := filepath.Join(config.OutDir, language.SubDir)

		templateLanguageRoot := filepath.Join(templateRoot, languageName.String(), "global")
		if err := walkTemplateDir(
			fileSystem,
			templateLanguageRoot,
			languageOutputPath,
			context,
		); err != nil {
			return err
		}

		templateLanguageModuleRoot := filepath.Join(templateRoot, languageName.String(), "module")
		for _, module := range config.Modules {
			context.Module = module

			fmt.Printf("\tGenerating templates for module: %s\n", module.Name)

			moduleOutputPath := filepath.Join(languageOutputPath, language.GetTemplateCompilePath(module))
			if err := walkTemplateDir(
				fileSystem,
				templateLanguageModuleRoot,
				moduleOutputPath,
				context,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func walkTemplateDir(fileSystem fs.FS, from string, to string, context *context) error {
	return fs.WalkDir(fileSystem, from, func(templateFilePath string, file os.DirEntry, err error) error {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}

		if err != nil {
			return err
		}

		if file.IsDir() {
			if templateFilePath == from {
				return nil
			}

			err = walkTemplateDir(fileSystem, templateFilePath, filepath.Join(to, file.Name()), context)
			if err != nil {
				return err
			}

			return nil
		}

		var template *templ.Template
		if template, err = template.ParseFS(fileSystem, templateFilePath); err != nil {
			return err
		}

		return createFileFromTemplate(template, to, file.Name(), context)
	})
}

func createFileFromTemplate(
	t *templ.Template,
	outputFilePath string,
	outputFileName string,
	context *context,
) error {
	err := os.MkdirAll(outputFilePath, os.ModePerm)
	if err != nil {
		return err
	}

	var file *os.File
	file, err = os.Create(filepath.Join(outputFilePath, outputFileName))
	if err != nil {
		return err
	}

	if err = t.Execute(file, context); err != nil {
		return err
	}

	err = file.Close()

	return err
}
