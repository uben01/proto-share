package template

import (
	. "config"
	"config/language"
	. "config/module"
)

type RenderConfig struct {
	Config   *Config
	Module   *Module
	Language *language.Config
}
