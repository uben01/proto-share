package template

import (
	"strings"
	templ "text/template"

	log "github.com/sirupsen/logrus"
)

const maxRecursionDepth = 10

func ProcessTemplateRecursively(templateString string, CTX interface{}) string {
	return processTemplateRecursively(templateString, CTX, 0)
}

func processTemplateRecursively(templateString string, CTX interface{}, depth int) string {
	if depth >= maxRecursionDepth {
		log.Panicf("max recursion depth (%d) exceeded", maxRecursionDepth)
	}

	template, err := templ.New("").Funcs(customFunctions).Parse(templateString)
	if err != nil {
		log.Panicf("failed to parse template: %s", err)
	}

	var buf strings.Builder
	if err = template.Execute(&buf, CTX); err != nil {
		log.Panicf("failed to execute template: %s", err)
	}

	if strings.Contains(buf.String(), "{{") {
		return processTemplateRecursively(buf.String(), CTX, depth+1)
	}

	return buf.String()
}
