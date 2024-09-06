package templating

import (
	"fmt"
	"strings"
	templ "text/template"

	. "github.com/uben01/proto-share/internal/context"
)

const maxRecursionDepth = 10

func ProcessTemplateRecursively(
	templateString string,
	depth int,
) (string, error) {
	if depth >= maxRecursionDepth {
		return "", fmt.Errorf("max recursion depth (%d) exceeded", maxRecursionDepth)
	}

	template, err := templ.New("").Funcs(customFunctions).Parse(templateString)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err = template.Execute(&buf, CTX); err != nil {
		return "", err
	}

	if strings.Contains(buf.String(), "{{") {
		return ProcessTemplateRecursively(buf.String(), depth+1)
	}

	return buf.String(), nil
}
