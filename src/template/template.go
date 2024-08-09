package template

import (
	"embed"
	"io/fs"
	. "os"
	"path/filepath"
	"text/template"

	. "proto-share/src/module"
	. "proto-share/src/param"
)

var templateRoot = filepath.Join("build", "templates")

func GenerateTemplates(embedFileSystem embed.FS, params *Param, modules []*Module) error {
	templateParam := TemplateParam{Param: params}

	for languageName, language := range params.Languages {
		templateParam.Language = &language

		languageOutputPath := filepath.Join(params.OutDir, language.SubDir)
		err := MkdirAll(languageOutputPath, ModePerm)
		if err != nil {
			return err
		}

		for _, module := range modules {
			err := MkdirAll(filepath.Join(languageOutputPath, language.ModulePath, module.Name), ModePerm)
			if err != nil {
				return err
			}
		}

		templateLanguageRoot := filepath.Join(templateRoot, string(languageName), "global")
		err = renderTemplates(embedFileSystem, templateLanguageRoot, languageOutputPath, templateParam)
		if err != nil {
			return err
		}

		templateLanguageModuleRoot := filepath.Join(templateRoot, string(languageName), "module")
		for _, module := range modules {
			templateParam.Module = module

			moduleOutputPath := filepath.Join(params.OutDir, language.SubDir, language.ModulePath, module.Name)
			err = renderTemplates(embedFileSystem, templateLanguageModuleRoot, moduleOutputPath, templateParam)

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func renderTemplates(embedFileSystem embed.FS, from string, to string, templateParam TemplateParam) error {
	return fs.WalkDir(embedFileSystem, from, func(path string, d DirEntry, err error) error {
		if d.IsDir() {
			if path == from {
				return nil
			}

			err = MkdirAll(filepath.Join(to, d.Name()), ModePerm)
			if err != nil {
				return err
			}
			err = renderTemplates(embedFileSystem, path, filepath.Join(to, d.Name()), templateParam)
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

		err = t.Execute(file, templateParam)
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
