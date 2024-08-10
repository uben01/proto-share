package config

import (
	"os"

	"gopkg.in/yaml.v2"

	. "proto-share/src/config/language"
	. "proto-share/src/config/module"
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

	var mergedLanguages = make(map[Language]LanguageConfig, len(config.Languages))
	for languageName, languageConfig := range config.Languages {
		lang, err := MergeWithDefault(languageName, languageConfig)
		if err != nil {
			return nil, err
		}
		mergedLanguages[languageName] = *lang
	}

	config.Languages = mergedLanguages

	return &config, err
}
