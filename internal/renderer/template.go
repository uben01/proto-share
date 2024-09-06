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
	"github.com/uben01/proto-share/internal/templating"
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
			CTX.Module = module

			fmt.Printf("\tGenerating templates for module: %s\n", module.Name)

			moduleOutputPath := filepath.Join(languageOutputPath, language.GetModuleTemplatePath(module))
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
		if len(fileContent) == 0 {
			return nil
		}

		dir := strings.TrimPrefix(filepath.Dir(templateFilePath), from)

		return createFileFromTemplate(
			string(fileContent),
			filepath.Join(to, dir),
			file.Name(),

			os.MkdirAll,
			os.Create,
		)
	})
}

var createFileFromTemplate = func(
	fileContent string,
	outputFilePath string,
	outputFileName string,

	mkdirAll func(path string, perm os.FileMode) error,
	createFile func(path string) (*os.File, error),
) error {
	err := mkdirAll(outputFilePath, os.ModePerm)
	if err != nil {
		return err
	}

	var file *os.File
	file, err = createFile(filepath.Join(outputFilePath, outputFileName))
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	processedTemplate, err := templating.ProcessTemplateRecursively(fileContent, 0)
	if err != nil {
		return err
	}

	_, err = file.WriteString(processedTemplate)

	return err
}
