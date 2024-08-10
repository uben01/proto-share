package config

import (
	. "proto-share/src/config/language"
	. "proto-share/src/config/module"
)

type TemplateConfig struct {
	Config   *Config
	Module   *Module
	Language *LanguageConfig
}
