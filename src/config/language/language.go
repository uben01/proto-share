package language

import (
	"fmt"
	"reflect"
)

type LanguageName string

const (
	Java LanguageName = "java"
)

type Language struct {
	SubDir               string            `yaml:"subDirName"`
	ModulePath           string            `yaml:"modulePath"`
	ProtoOutputDir       string            `yaml:"protoOutputDir"`
	ProtocCommand        string            `yaml:"protocCommand"`
	AdditionalParameters map[string]string `yaml:"additionalParameters"`
}

var defaultMapping = map[LanguageName]*Language{
	Java: defaultJava(),
}

func MergeWithDefault(languageName LanguageName, actualLanguage *Language) (*Language, error) {
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
