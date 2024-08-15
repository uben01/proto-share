package language

func defaultPHP() *Language {
	return &Language{
		SubDir:            "php",
		ModulePath:        ".",
		SeparateModuleDir: false,
		ProtoOutputDir:    ".",
		ProtocCommand:     "php_out",
	}
}
