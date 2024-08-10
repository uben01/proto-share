package config

import (
	. "proto-share/src/language"
	. "proto-share/src/module"
)

type TemplateConfig struct {
	Config   *Config
	Module   *Module
	Language *LanguageConfig
}
