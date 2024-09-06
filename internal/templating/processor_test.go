package templating

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

	output, err := ProcessTemplateRecursively(templateString, ctx)
	assert.Nil(t, err)
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

	output, err := ProcessTemplateRecursively(templateString, ctx)
	assert.Nil(t, err)
	assert.Equal(t, "module-name", output)
}

func TestProcessTemplateRecursively_InfiniteRecursion_ReturnsAfterMaxDepthReached(t *testing.T) {
	templateString := `{{ .Module.Path }}`
	ctx := &ctx{
		Module: &module{
			Path: "{{ .Module.Path }}",
		},
	}

	_, err := ProcessTemplateRecursively(templateString, ctx)
	assert.Error(t, err)
	assert.Equal(t, "max recursion depth (10) exceeded", err.Error())
}
