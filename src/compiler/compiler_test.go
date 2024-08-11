package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "config"
	"config/language"
	. "config/module"
)

func TestCompileModules_WithoutModules(t *testing.T) {
	config := &Config{
		Modules: []*Module{},
	}

	err := CompileModules(config)

	assert.NotNil(t, err)
	assert.Errorf(t, err, "no modules defined")
}

func TestCompileModules_WithoutLanguages(t *testing.T) {
	config := &Config{
		Modules: []*Module{
			{},
		},
		Languages: map[language.Name]*language.Config{},
	}

	err := CompileModules(config)

	assert.NotNil(t, err)
	assert.Errorf(t, err, "no languages defined")
}
