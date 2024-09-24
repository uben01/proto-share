package language

func defaultPHP() *Language {
	return &Language{
		ModuleTemplatePath: "{{ .Module.Name | snakeCase }}",
		ModuleCompilePath:  "{{ .Module.Name | snakeCase }}",
		SubDir:             "php",
		ProtocCommand:      "php_out",
		AdditionalParameters: map[string]string{
			"phpVersion": "8.1",
			"vendor":     "",
			"moduleName": "{{ .Module.Name | kebabCase }}",
		},
	}
}
