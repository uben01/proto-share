package render

import (
	"errors"
	"fmt"
	"io/fs"
	. "os"
	"path/filepath"
	"text/template"

	. "config"
)

var templateRoot = filepath.Join("build", "templates")

func GenerateTemplates(fileSystem fs.FS, config *Config) error {
	context := &context{Config: config}

	for languageName, language := range config.Languages {
		context.Language = language

		languageOutputPath := filepath.Join(config.OutDir, language.SubDir)

		for _, module := range config.Modules {
			if err := MkdirAll(
				filepath.Join(languageOutputPath, language.GetModulePath(module)),
				ModePerm,
			); err != nil {
				return err
			}
		}

		fmt.Printf("Generating templates for language: %s\n", languageName)

		templateLanguageRoot := filepath.Join(templateRoot, languageName.String(), "global")
		if err := renderTemplates(
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

			moduleOutputPath := filepath.Join(languageOutputPath, language.GetModulePath(module))

			if err := renderTemplates(
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

func renderTemplates(fileSystem fs.FS, from string, to string, context *context) error {
	return fs.WalkDir(fileSystem, from, func(path string, d DirEntry, err error) error {
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

			err = MkdirAll(filepath.Join(to, d.Name()), ModePerm)
			if err != nil {
				return err
			}
			err = renderTemplates(fileSystem, path, filepath.Join(to, d.Name()), context)
			if err != nil {
				return err
			}

			return nil
		}

		outputPath := filepath.Join(to, d.Name())

		t, err := template.ParseFS(fileSystem, path)
		if err != nil {
			return err
		}

		file, err := Create(outputPath)
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
