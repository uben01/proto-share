package template

import (
	"reflect"
	"text/template"

	"github.com/ettle/strcase"
)

var customFunctions = template.FuncMap{
	"required":   required,
	"kebabCase":  strcase.ToKebab,
	"snakeCase":  strcase.ToSnake,
	"camelCase":  strcase.ToCamel,
	"pascalCase": strcase.ToPascal,
}

func required(v interface{}) interface{} {
	panicMessage := "Required field not provided. Check your configuration!"

	if v == nil {
		panic(panicMessage)
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.String, reflect.Map, reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			panic(panicMessage)
		}
	case reflect.Ptr:
		if val.IsNil() {
			panic(panicMessage)
		}
	default:
		panic("Unhandled case. Please report this issue!")
	}
	return v
}
