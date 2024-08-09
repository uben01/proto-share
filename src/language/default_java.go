package language

func DefaultJava() LanguageParams {
	return LanguageParams{
		Name:           Java,
		SubDir:         "java",
		ModulePath:     ".",
		ProtoOutputDir: "src/main/java",
		ProtocCommand:  "java_out",
		AdditionalParams: map[string]string{
			"version": "21",
			"jarPath": "${rootDir}/build/libs",
		},
	}
}
