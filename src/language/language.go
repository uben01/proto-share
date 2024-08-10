package language

import "reflect"

type Language string

const (
	Java Language = "java"
)

type LanguageConfig struct {
	SubDir               string            `yaml:"subDirName"`
	ModulePath           string            `yaml:"modulePath"`
	ProtoOutputDir       string            `yaml:"protoOutputDir"`
	ProtocCommand        string            `yaml:"protocCommand"`
	AdditionalParameters map[string]string `yaml:"additionalParameters"`
}

func MergeWithDefault(
	actualLanguageConfig LanguageConfig,
	defaultLanguageConfig LanguageConfig,
) LanguageConfig {
	merged := actualLanguageConfig
	defaultVal := reflect.ValueOf(defaultLanguageConfig)
	mergedVal := reflect.ValueOf(&merged).Elem()

	for i := 0; i < defaultVal.NumField(); i++ {
		field := defaultVal.Type().Field(i)
		defaultField := defaultVal.Field(i)
		mergedField := mergedVal.FieldByName(field.Name)

		if isEmptyValue(mergedField) {
			mergedField.Set(defaultField)
		}
	}

	return merged
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
