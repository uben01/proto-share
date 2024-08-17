package language

func defaultTS() *Language {
	return &Language{
		ModulePathTemplate: "",
		SubDir:             "ts",
		ProtocCommand:      "ts_out",
	}
}
