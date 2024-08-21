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
	ModulePathTemplate   string            `yaml:"modulePathTemplate"`
	SubDir               string            `yaml:"subDirName"`
	ProtocCommand        string            `yaml:"protocCommand"`
	AdditionalParameters map[string]string `yaml:"additionalParameters"`
}

func (language Language) GetModulePath(module *module.Module) string {
	return strings.ReplaceAll(language.ModulePathTemplate, "{module}", module.Name)
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
	case reflect.Bool:
		return !v.Bool()
	default:
		panic("unhandled default case")
	}
	return false
}
