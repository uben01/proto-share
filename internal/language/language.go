package language

import (
	"fmt"
	"strings"

	"dario.cat/mergo"

	"github.com/uben01/proto-share/internal/module"
)

type Name string

const (
	Java       Name = "Java"
	PHP        Name = "PHP"
	TypeScript Name = "TypeScript"
)

func (name Name) String() string {
	return strings.ToLower(string(name))
}

type Language struct {

	// protoc output dir:
	// {config.outDir}/{language.subDir}/{language.moduleCompilePath}
	ModuleCompilePath string `yaml:"moduleCompilePath"`

	// templates are copied
	// from `assets/templates/{language}/module`
	// to 	`{config.outDir}/{language.subDir}/{language.moduleTemplatePath}`
	ModuleTemplatePath string `yaml:"moduleTemplatePath"`

	// output subdirectory name for language
	SubDir string `yaml:"subDirName"`

	// protoc command to generate code for language e.g. `java_out`, `php_out`...
	ProtocCommand string `yaml:"protocCommand"`

	// generate publish script for language
	EnablePublish *bool `yaml:"enablePublish"`

	// additional parameters to be passed for templating
	// documented for every language, or can be used for custom templates
	AdditionalParameters map[string]string `yaml:"additionalParameters"`
}

func (language Language) GetModuleCompilePath(module *module.Module) string {
	return strings.ReplaceAll(language.ModuleCompilePath, "{module}", module.Name)
}

func (language Language) GetModuleTemplatePath(module *module.Module) string {
	return strings.ReplaceAll(language.ModuleTemplatePath, "{module}", module.Name)
}

var defaultMapping = map[Name]*Language{
	Java:       defaultJava(),
	PHP:        defaultPHP(),
	TypeScript: defaultTS(),
}

func MergeWithDefault(languageName Name, actualLanguage *Language) (*Language, error) {
	defaultLanguageConfig := defaultMapping[languageName]

	if defaultLanguageConfig == nil {
		return nil, fmt.Errorf("unsupported language: %s", languageName)
	}

	if actualLanguage == nil {
		return defaultMapping[languageName], nil
	}

	var merged = &Language{}
	var err error
	if err = mergo.Merge(merged, defaultLanguageConfig); err != nil {
		return nil, err
	}

	err = mergo.Merge(merged, *actualLanguage, mergo.WithOverride)

	return merged, err
}
