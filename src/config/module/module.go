package module

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Module struct {
	Name    string `yaml:"name"`
	Hash    string `yaml:"hash"`
	Version int    `yaml:"version"`
	Path    string `yaml:"path"`
}

func GetAllModules(root string) ([]*Module, error) {
	files, err := findModuleYmlFiles(root)
	if err != nil {
		return nil, err
	}

	var modules []*Module
	for _, file := range files {
		module, err := readAndParseModuleYml(file)
		if err != nil {
			return nil, err
		}
		modules = append(modules, module)
	}

	return modules, nil
}

func findModuleYmlFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "module.yml" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func readAndParseModuleYml(filePath string) (*Module, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var module Module
	err = yaml.Unmarshal(data, &module)
	if err != nil {
		return nil, err
	}

	return &module, nil
}
