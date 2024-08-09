package language

type Language struct {
	Name             string            `yaml:"name"`
	SubDir           string            `yaml:"subDirName"`
	ModulePath       string            `yaml:"modulePath"`
	ProtoOutputDir   string            `yaml:"protoOutputDir"`
	ProtocCommand    string            `yaml:"protocCommand"`
	AdditionalParams map[string]string `yaml:"additionalParams"`
}
