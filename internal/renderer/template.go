package renderer

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

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
		if err := createTemplateFiles(
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
			if err := createTemplateFiles(
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

func createTemplateFiles(fileSystem fs.FS, from string, to string, context *context) error {
	return fs.WalkDir(fileSystem, from, func(path string, d os.DirEntry, err error) error {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}

		if err != nil {
			return err
		}

		if d.IsDir() {
			if path == from {
				return nil
			}

			err = createTemplateFiles(fileSystem, path, filepath.Join(to, d.Name()), context)
			if err != nil {
				return err
			}

			return nil
		}

		err = os.MkdirAll(to, os.ModePerm)
		if err != nil {
			return err
		}

		t, err := template.ParseFS(fileSystem, path)
		if err != nil {
			return err
		}

		file, err := os.Create(filepath.Join(to, d.Name()))
		if err != nil {
			return err
		}

		err = t.Execute(file, context)
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}

		return nil
	})
}
