package language

func defaultTS() *Language {
	return &Language{
		SubDir:            "ts",
		ModulePath:        ".",
		SeparateModuleDir: false,
		ProtoOutputDir:    ".",
		ProtocCommand:     "ts_out",
	}
}
