package param

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	. "proto-share/src/language"
	. "proto-share/src/module"
)

type Param struct {
	ProjectName string                     `yaml:"projectName"`
	InDir       string                     `yaml:"inDir"`
	OutDir      string                     `yaml:"outDir"`
	Languages   map[Language]LanguageParam `yaml:"languages"`

	Modules []*Module
}

func ParseParams(configPath string) (*Param, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var params Param
	err = yaml.Unmarshal(data, &params)

	var languages = make(map[Language]LanguageParam, len(params.Languages))
	for languageName, language := range params.Languages {
		switch languageName {
		case Java:
			lang := MergeWithDefault(language, DefaultJava())
			languages[languageName] = lang
		default:
			return nil, fmt.Errorf("unsupported language: %s", languageName)
		}
	}

	params.Languages = languages

	return &params, err
}
