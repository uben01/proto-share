package languages

func Java() Language {
	java := Language{
		subDir:         "java",
		modulePath:     ".",
		protoOutputDir: "src/main/java",
		protocCommand:  "java_out",
	}

	return java
}
