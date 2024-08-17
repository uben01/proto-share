package language

func defaultPHP() *Language {
	return &Language{
		ModulePathTemplate: "",
		SubDir:             "php",
		ProtocCommand:      "php_out",
	}
}
