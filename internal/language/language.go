package language

import (
	"fmt"
	"reflect"
	"strings"

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
	// to 	 `{config.outDir}/{language.subDir}/{language.moduleTemplatePath}`
	ModuleTemplatePath string `yaml:"moduleTemplatePath"`

	// output subdirectory name for language
	SubDir string `yaml:"subDirName"`

	// protoc command to generate code for language e.g. `java_out`, `php_out`...
	ProtocCommand string `yaml:"protocCommand"`

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

	var merged *Language
	if actualLanguage != nil {
		merged = actualLanguage
	} else {
		merged = &Language{}
	}

	defaultVal := reflect.ValueOf(*defaultLanguageConfig)
	mergedVal := reflect.ValueOf(merged).Elem()

	for i := 0; i < defaultVal.NumField(); i++ {
		field := defaultVal.Type().Field(i)
		defaultField := defaultVal.Field(i)
		mergedField := mergedVal.FieldByName(field.Name)

		if isEmptyValue(mergedField) {
			mergedField.Set(defaultField)
		}
	}

	return merged, nil
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String, reflect.Map, reflect.Slice, reflect.Array:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		panic("unhandled default case")
	}
	return false
}
