package main

import (
	"os"
	"path/filepath"
	. "proto-share/src/languages"
	"text/template"
)

func renderTemplates(from string, to string, templateData TemplateData) error {
	return filepath.WalkDir(from, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			if path == from {
				return nil
			}

			err = os.MkdirAll(filepath.Join(to, d.Name()), os.ModePerm)
			if err != nil {
				return err
			}
			err = renderTemplates(path, filepath.Join(to, d.Name()), templateData)
			if err != nil {
				return err
			}

			return nil
		}

		outputPath := filepath.Join(to, d.Name())
		t := template.Must(template.ParseFiles(path))
		if err != nil {
			return err
		}

		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}

		err = t.Execute(file, templateData)
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

func generateTemplates(languages []Language, modules []*Module, outDir string, templateData TemplateData) error {
	for _, language := range languages {
		languageOutputPath := filepath.Join(outDir, language.SubDir())
		err := os.MkdirAll(languageOutputPath, os.ModePerm)
		if err != nil {
			return err
		}

		for _, module := range modules {
			err := os.MkdirAll(filepath.Join(languageOutputPath, language.ModulePath(), module.Name), os.ModePerm)
			if err != nil {
				return err
			}
		}

		templateLanguageRoot := filepath.Join("templates", language.SubDir(), "global")
		err = renderTemplates(templateLanguageRoot, languageOutputPath, templateData)
		if err != nil {
			return err
		}

		templateLanguageModuleRoot := filepath.Join("templates", language.SubDir(), "module")
		for _, module := range modules {
			templateData.CurrentModule = module
			moduleOutputPath := filepath.Join(outDir, language.SubDir(), language.ModulePath(), module.Name)
			err = renderTemplates(templateLanguageModuleRoot, moduleOutputPath, templateData)

			if err != nil {
				return err
			}
		}
	}
	return nil
}
