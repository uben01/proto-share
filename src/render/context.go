package render

import (
	. "config"
	. "config/language"
	. "config/module"
)

type Context struct {
	Config   *Config
	Module   *Module
	Language *Language
}
