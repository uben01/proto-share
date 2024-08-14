package language

func defaultTS() *Language {
	return &Language{
		SubDir:         "ts",
		ModulePath:     ".",
		ProtoOutputDir: ".",
		ProtocCommand:  "ts_out",
	}
}
