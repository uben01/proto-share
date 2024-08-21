package render

import (
	. "github.com/uben01/proto-share/internal/config"
	. "github.com/uben01/proto-share/internal/language"
	. "github.com/uben01/proto-share/internal/module"
)

type context struct {
	Config   *Config
	Module   *Module
	Language *Language
}
