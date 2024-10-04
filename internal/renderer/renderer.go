package renderer

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	. "github.com/uben01/proto-share/internal/config"
	. "github.com/uben01/proto-share/internal/context"
	"github.com/uben01/proto-share/internal/template"
)

var templateRoot = filepath.Join("assets", "templates")

func RenderTemplates(fileSystem fs.FS, config *Config) {
	if len(config.Languages) == 0 {
		log.Panicf("no languages found in config")
	}

	if len(config.Modules) == 0 {
		log.Panicf("no modules found in config")
	}

	for languageName, language := range config.Languages {
		CTX.Language = language

		log.Infof("Generating templates for language: %s", languageName)

		languageOutputPath := filepath.Join(config.OutDir, language.SubDir)

		templateLanguageRoot := filepath.Join(templateRoot, languageName.String(), "global")
		walkTemplateDir(fileSystem, templateLanguageRoot, languageOutputPath)

		templateLanguageModuleRoot := filepath.Join(templateRoot, languageName.String(), "module")
		for _, module := range config.Modules {
			if !module.Changed {
				log.Debugf("Module %s has not been changed", module.Name)

				if !config.ForceGeneration {
					continue
				}
			}

			CTX.Module = module
			log.Infof("Generating templates for module: %s", module.Name)

			moduleOutputPath := filepath.Join(
				languageOutputPath,
				template.ProcessTemplateRecursively(language.ModuleTemplatePath, CTX),
			)
			walkTemplateDir(fileSystem, templateLanguageModuleRoot, moduleOutputPath)
		}
	}
}

var walkTemplateDir = func(
	fileSystem fs.FS,
	from string,
	to string,
) {
	err := fs.WalkDir(fileSystem, from, func(templateFilePath string, file os.DirEntry, err error) error {
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
			log.Panicf("error reading template file: %s", err)
		}

		processedTemplate := template.ProcessTemplateRecursively(string(fileContent), CTX)

		dir := strings.TrimPrefix(filepath.Dir(templateFilePath), from)

		return createFileFromTemplate(processedTemplate, filepath.Join(to, dir), file.Name())
	})
	if err != nil {
		log.Panicf("error walking template directory: %s", err)
	}
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
