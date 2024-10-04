package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ctx struct {
	Module *module
}

type module struct {
	Name string
	Path string
}

func TestProcessTemplateRecursively_OneLevel_OutputReturned(t *testing.T) {
	templateString := `{{ .Module.Path }}`
	ctx := &ctx{
		Module: &module{
			Path: "module",
		},
	}

	output := ProcessTemplateRecursively(templateString, ctx)
	assert.Equal(t, "module", output)
}

func TestProcessTemplateRecursively_TwoLevelsWithCustomFunctions_OutputReturned(t *testing.T) {
	templateString := `{{ .Module.Path }}`
	ctx := &ctx{
		Module: &module{
			Name: "moduleName",
			Path: "{{ .Module.Name | kebabCase }}",
		},
	}

	output := ProcessTemplateRecursively(templateString, ctx)
	assert.Equal(t, "module-name", output)
}

func TestProcessTemplateRecursively_InfiniteRecursion_Panics(t *testing.T) {
	templateString := `{{ .Module.Path }}`
	ctx := &ctx{
		Module: &module{
			Path: "{{ .Module.Path }}",
		},
	}

	assert.Panics(t, func() {
		_ = ProcessTemplateRecursively(templateString, ctx)
	})
}
