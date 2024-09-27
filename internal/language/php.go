package language

func defaultPHP() *Language {
	return &Language{
		ModuleTemplatePath: "{{ .Module.Name | snakeCase }}",
		ModuleCompilePath:  "{{ .Module.Name | snakeCase }}/src",
		SubDir:             "php",
		ProtocCommand:      "php_out",
		AdditionalParameters: map[string]string{
			"phpVersion": "^8",
			"vendor":     "",
			"moduleName": "{{ .Module.Name | kebabCase }}",
		},
	}
}
