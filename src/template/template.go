package template

import (
	"embed"
	"io/fs"
	. "os"
	"path/filepath"
	"text/template"
)

var templateRoot = filepath.Join("build", "templates")

func GenerateTemplates(embedFileSystem embed.FS, renderConfig *RenderConfig) error {
	for languageName, language := range renderConfig.Config.Languages {
		renderConfig.Language = language

		languageOutputPath := filepath.Join(renderConfig.Config.OutDir, language.SubDir)
		if err := MkdirAll(languageOutputPath, ModePerm); err != nil {
			return err
		}

		for _, module := range renderConfig.Config.Modules {
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
			renderConfig,
		); err != nil {
			return err
		}

		templateLanguageModuleRoot := filepath.Join(templateRoot, string(languageName), "module")
		for _, module := range renderConfig.Config.Modules {
			renderConfig.Module = module

			moduleOutputPath := filepath.Join(renderConfig.Config.OutDir, language.SubDir, language.ModulePath, module.Name)
			if err := renderTemplates(
				embedFileSystem,
				templateLanguageModuleRoot,
				moduleOutputPath,
				renderConfig,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func renderTemplates(embedFileSystem embed.FS, from string, to string, renderConfig *RenderConfig) error {
	return fs.WalkDir(embedFileSystem, from, func(path string, d DirEntry, err error) error {
		if d.IsDir() {
			if path == from {
				return nil
			}

			err = MkdirAll(filepath.Join(to, d.Name()), ModePerm)
			if err != nil {
				return err
			}
			err = renderTemplates(embedFileSystem, path, filepath.Join(to, d.Name()), renderConfig)
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

		err = t.Execute(file, renderConfig)
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
