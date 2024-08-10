package template

import (
	"embed"
	"io/fs"
	. "os"
	"path/filepath"
	"text/template"

	. "proto-share/src/config"
)

var templateRoot = filepath.Join("build", "templates")

func GenerateTemplates(embedFileSystem embed.FS, config *Config) error {
	templateConfig := TemplateConfig{Config: config}

	for languageName, language := range config.Languages {
		templateConfig.Language = language

		languageOutputPath := filepath.Join(config.OutDir, language.SubDir)
		if err := MkdirAll(languageOutputPath, ModePerm); err != nil {
			return err
		}

		for _, module := range config.Modules {
			if err := MkdirAll(
				filepath.Join(languageOutputPath, language.ModulePath, module.Name),
				ModePerm,
			); err != nil {
				return err
			}
		}

		templateLanguageRoot := filepath.Join(templateRoot, string(languageName), "global")
		if err := renderTemplates(
			embedFileSystem,
			templateLanguageRoot,
			languageOutputPath,
			templateConfig,
		); err != nil {
			return err
		}

		templateLanguageModuleRoot := filepath.Join(templateRoot, string(languageName), "module")
		for _, module := range config.Modules {
			templateConfig.Module = module

			moduleOutputPath := filepath.Join(config.OutDir, language.SubDir, language.ModulePath, module.Name)
			if err := renderTemplates(
				embedFileSystem,
				templateLanguageModuleRoot,
				moduleOutputPath,
				templateConfig,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func renderTemplates(embedFileSystem embed.FS, from string, to string, templateConfig TemplateConfig) error {
	return fs.WalkDir(embedFileSystem, from, func(path string, d DirEntry, err error) error {
		if d.IsDir() {
			if path == from {
				return nil
			}

			err = MkdirAll(filepath.Join(to, d.Name()), ModePerm)
			if err != nil {
				return err
			}
			err = renderTemplates(embedFileSystem, path, filepath.Join(to, d.Name()), templateConfig)
			if err != nil {
				return err
			}

			return nil
		}

		outputPath := filepath.Join(to, d.Name())

		t := template.Must(template.ParseFS(embedFileSystem, path))
		if err != nil {
			return err
		}

		file, err := Create(outputPath)
		if err != nil {
			return err
		}

		err = t.Execute(file, templateConfig)
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
