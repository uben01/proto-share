package context

import (
	. "github.com/uben01/proto-share/internal/config"
	. "github.com/uben01/proto-share/internal/language"
	. "github.com/uben01/proto-share/internal/module"
)

type Context struct {
	Config   *Config
	Module   *Module
	Language *Language
	Env      Environment
}

var CTX = func() *Context {
	env := prepareEnvironment()

	return &Context{Env: env}
}()
