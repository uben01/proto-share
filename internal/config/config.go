package config

import (
	"fmt"
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"

	. "github.com/uben01/proto-share/internal/language"
	. "github.com/uben01/proto-share/internal/module"
)

type Config struct {
	// project name for the build tools
	ProjectName string `yaml:"projectName"`

	// input directory for the proto files
	InDir string `yaml:"inDir"`

	// output directory for the generated code
	OutDir string `yaml:"outDir"`

	// languages to generate code for
	Languages map[Name]*Language `yaml:"languages"`

	Modules []*Module
}

var fileSystem = os.DirFS(".")

var mergeWithDefault = MergeWithDefault

func ParseConfig(configPath string) (*Config, error) {
	data, err := fs.ReadFile(fileSystem, configPath)
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
		lang, err := mergeWithDefault(languageName, languageConfig)
		if err != nil {
			return err
		}
		mergedLanguages[languageName] = lang
	}
	config.Languages = mergedLanguages

	return nil
}
