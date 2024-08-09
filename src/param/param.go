package param

import (
	"os"

	"gopkg.in/yaml.v2"

	. "proto-share/src/language"
	. "proto-share/src/module"
)

type Param struct {
	ProjectName string     `yaml:"projectName"`
	InDir       string     `yaml:"inDir"`
	OutDir      string     `yaml:"outDir"`
	Languages   []Language `yaml:"languages"`

	Module  *Module
	Modules []*Module
}

func ParseParams(configPath string) (*Param, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var params Param
	err = yaml.Unmarshal(data, &params)

	return &params, err
}
