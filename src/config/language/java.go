package language

func defaultJava() *Config {
	return &Config{
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
