package param

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	. "proto-share/src/language"
	. "proto-share/src/module"
)

type Param struct {
	ProjectName string           `yaml:"projectName"`
	InDir       string           `yaml:"inDir"`
	OutDir      string           `yaml:"outDir"`
	Languages   []LanguageParams `yaml:"languages"`

	Module  *Module
	Modules []*Module
}

func ParseParams(configPath string) (*Param, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var params Param
	err = yaml.Unmarshal(data, &params)

	var languages []LanguageParams
	for _, language := range params.Languages {
		switch language.Name {
		case Java:
			lang := MergeWithDefault(language, DefaultJava())
			languages = append(languages, lang)
		default:
			return nil, fmt.Errorf("unsupported language: %s", language.Name)
		}
	}

	params.Languages = languages

	return &params, err
}
