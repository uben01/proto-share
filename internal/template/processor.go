package template

import (
	"fmt"
	"strings"
	templ "text/template"
)

const maxRecursionDepth = 10

func Must(res string, err error) string {
	if err != nil {
		panic(err)
	}

	return res
}

func ProcessTemplateRecursively(templateString string, CTX interface{}) (string, error) {
	return processTemplateRecursively(templateString, CTX, 0)
}

func processTemplateRecursively(templateString string, CTX interface{}, depth int) (string, error) {
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
		return processTemplateRecursively(buf.String(), CTX, depth+1)
	}

	return buf.String(), nil
}
