package config

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	log "github.com/sirupsen/logrus"

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

	// force generation of code even if no changes are detected
	ForceGeneration bool `yaml:"forceGeneration"`

	// languages to generate code for
	Languages map[Name]*Language `yaml:"languages"`

	Modules []*Module
}

func (c *Config) AnyModuleChanged() bool {
	for _, module := range c.Modules {
		if module.Changed {
			return true
		}
	}

	return false
}

var fileSystem = os.DirFS(".")

var mergeWithDefault = MergeWithDefault

func ParseConfig(configPath string) *Config {
	data, err := fs.ReadFile(fileSystem, filepath.Clean(configPath))
	if err != nil {
		log.Panicf("Error reading config file: %s: %v", configPath, err)
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		log.Panicf("Failed to parse config file: %s", err)
	}

	mergeConfigWithDefaults(&config)

	log.Debugf("Parsed config at: %s", configPath)

	return &config
}

func mergeConfigWithDefaults(config *Config) {
	var mergedLanguages = make(map[Name]*Language, len(config.Languages))
	for languageName, languageConfig := range config.Languages {
		lang := mergeWithDefault(languageName, languageConfig)
		mergedLanguages[languageName] = lang
	}
	config.Languages = mergedLanguages
}
