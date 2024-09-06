package templating

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/uben01/proto-share/internal/context"
	"github.com/uben01/proto-share/internal/module"
)

func TestProcessTemplateRecursively_OneLevel_OutputReturned(t *testing.T) {
	templateString := `{{ .Module.Path }}`
	CTX = &Context{
		Module: &module.Module{
			Path: "module",
		},
	}

	output, err := ProcessTemplateRecursively(templateString, 0)
	assert.Nil(t, err)
	assert.Equal(t, "module", output)
}

func TestProcessTemplateRecursively_TwoLevelsWithCustomFunctions_OutputReturned(t *testing.T) {
	templateString := `{{ .Module.Path }}`
	CTX = &Context{
		Module: &module.Module{
			Name: "moduleName",
			Path: "{{ .Module.Name | kebabCase }}",
		},
	}

	output, err := ProcessTemplateRecursively(templateString, 0)
	assert.Nil(t, err)
	assert.Equal(t, "module-name", output)
}

func TestProcessTemplateRecursively_InfiniteRecursion_ReturnsAfterMaxDepthReached(t *testing.T) {
	templateString := `{{ .Module.Path }}`
	CTX = &Context{
		Module: &module.Module{
			Path: "{{ .Module.Path }}",
		},
	}

	_, err := ProcessTemplateRecursively(templateString, 0)
	assert.Error(t, err)
	assert.Equal(t, "max recursion depth (10) exceeded", err.Error())
}
