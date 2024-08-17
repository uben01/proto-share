package language

func defaultJava() *Language {
	return &Language{
		ModulePathTemplate: "{module}/src/main/java",
		SubDir:             "java",
		ProtocCommand:      "java_out",
		AdditionalParameters: map[string]string{
			"version": "21",
			"jarPath": "${rootDir}/build/libs",
		},
	}
}
