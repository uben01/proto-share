package language

func defaultTS() *Language {
	return &Language{
		ModuleTemplatePath: "",
		ModuleCompilePath:  "",
		SubDir:             "ts",
		ProtocCommand:      "ts_out",
	}
}
