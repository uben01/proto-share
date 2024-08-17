package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	. "config/language"
	. "config/module"
)

type Config struct {
	ProjectName string             `yaml:"projectName"`
	InDir       string             `yaml:"inDir"`
	OutDir      string             `yaml:"outDir"`
	Languages   map[Name]*Language `yaml:"languages"`

	Modules []*Module
}

func ParseConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err = mergeConfigWithDefaults(&config); err != nil {
		return nil, err
	}

	fmt.Printf("Parsed config at: %s\n", configPath)

	return &config, nil
}

func mergeConfigWithDefaults(config *Config) error {
	var mergedLanguages = make(map[Name]*Language, len(config.Languages))
	for languageName, languageConfig := range config.Languages {
		lang, err := MergeWithDefault(languageName, languageConfig)
		if err != nil {
			return err
		}
		mergedLanguages[languageName] = lang
	}
	config.Languages = mergedLanguages

	return nil
}
