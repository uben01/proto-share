package language

import (
	"strings"

	"dario.cat/mergo"
	log "github.com/sirupsen/logrus"
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

var defaultMapping = map[Name]*Language{
	Java:       defaultJava(),
	PHP:        defaultPHP(),
	TypeScript: defaultTS(),
}

func MergeWithDefault(languageName Name, actualLanguage *Language) *Language {
	defaultLanguageConfig := defaultMapping[languageName]

	if defaultLanguageConfig == nil {
		log.Panicf("unsupported language: %s", languageName)
	}

	if actualLanguage == nil {
		return defaultMapping[languageName]
	}

	var merged = &Language{}
	var err error
	if err = mergo.Merge(merged, defaultLanguageConfig); err != nil {
		log.Panicf("failed to merge default language config: %s", err)
	}

	err = mergo.Merge(merged, *actualLanguage, mergo.WithOverride)
	if err != nil {
		log.Panicf("failed to merge actual language config: %s", err)
	}

	return merged
}
