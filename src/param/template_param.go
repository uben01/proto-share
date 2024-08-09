package param

import (
	. "proto-share/src/language"
	. "proto-share/src/module"
)

type TemplateParam struct {
	Param    *Param
	Module   *Module
	Language *LanguageParam
}
