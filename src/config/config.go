package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	. "proto-share/src/language"
	. "proto-share/src/module"
)

type Config struct {
	ProjectName string                      `yaml:"projectName"`
	InDir       string                      `yaml:"inDir"`
	OutDir      string                      `yaml:"outDir"`
	Languages   map[Language]LanguageConfig `yaml:"languages"`

	Modules []*Module
}

func ParseConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)

	var languages = make(map[Language]LanguageConfig, len(config.Languages))
	for languageName, language := range config.Languages {
		switch languageName {
		case Java:
			lang := MergeWithDefault(language, DefaultJava())
			languages[languageName] = lang
		default:
			return nil, fmt.Errorf("unsupported language: %s", languageName)
		}
	}

	config.Languages = languages

	return &config, err
}
