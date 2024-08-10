package language

func defaultJava() *LanguageConfig {
	return &LanguageConfig{
		SubDir:         "java",
		ModulePath:     ".",
		ProtoOutputDir: "src/main/java",
		ProtocCommand:  "java_out",
		AdditionalParameters: map[string]string{
			"version": "21",
			"jarPath": "${rootDir}/build/libs",
		},
	}
}
