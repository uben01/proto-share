package language

func defaultPHP() *Language {
	return &Language{
		ModuleTemplatePath: "",
		ModuleCompilePath:  "",
		SubDir:             "php",
		ProtocCommand:      "php_out",
	}
}
