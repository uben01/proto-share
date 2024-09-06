package language

func defaultJava() *Language {
	return &Language{
		ModuleCompilePath:  "{module}/src/main/java",
		ModuleTemplatePath: "{module}",
		SubDir:             "java",
		ProtocCommand:      "java_out",
		AdditionalParameters: map[string]string{
			"version":            "21",
			"jarPath":            "${rootDir}/build/libs",
			"groupId":            "",
			"artifactId":         "",
			"repositoryUrl":      "",
			"repositoryUsername": "",
			"repositoryPassword": "",
		},
	}
}
