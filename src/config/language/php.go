package language

func defaultPHP() *Language {
	return &Language{
		SubDir:         "php",
		ModulePath:     ".",
		ProtoOutputDir: ".",
		ProtocCommand:  "php_out",
	}
}
