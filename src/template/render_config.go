package template

import (
	. "proto-share/src/config"
	"proto-share/src/config/language"
	. "proto-share/src/config/module"
)

type RenderConfig struct {
	Config   *Config
	Module   *Module
	Language *language.Config
}
