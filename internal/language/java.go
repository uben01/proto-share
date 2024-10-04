package language

func defaultJava() *Language {
	return &Language{
		ModuleCompilePath:  "{{ .Module.Name }}/src/main/java",
		ModuleTemplatePath: "{{ .Module.Name }}",
		SubDir:             "java",
		ProtocCommand:      "java_out",
		AdditionalParameters: map[string]string{
			"version":             "21",
			"jarPath":             "${rootDir}/build/libs",
			"enableMavenPublish:": "false",
			"groupId":             "",
			"artifactId":          "{{ .Module.Name | kebabCase }}",
			"repositoryUrl":       "",
			"repositoryUsername":  "",
			"repositoryPassword":  "",
		},
	}
}
