package languages

type Language struct {
	subDir         string // Subdirectory name in outDir
	modulePath     string // Module path relative to language subdirectory
	protoOutputDir string // Output directory for generated proto files relative to the modulePath
	protocCommand  string // Protoc argument to generate code for this language (e.g. "java_out")
}

func (l Language) SubDir() string {
	return l.subDir
}

func (l Language) ModulePath() string {
	return l.modulePath
}

func (l Language) ProtoOutputDir() string {
	return l.protoOutputDir
}

func (l Language) ProtocCommand() string {
	return l.protocCommand
}
