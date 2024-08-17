package render

import (
	. "config"
	. "config/language"
	. "config/module"
)

type context struct {
	Config   *Config
	Module   *Module
	Language *Language
}
