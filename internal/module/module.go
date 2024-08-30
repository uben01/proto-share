package module

import (
	"fmt"
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"
)

type Module struct {
	Name    string `yaml:"name"`
	Hash    string `yaml:"hash"`
	Version int    `yaml:"version"`
	Path    string `yaml:"path"`
}

var moduleFileName = "module.yml"

var fileSystem = os.DirFS(".")

func DiscoverModules(root string) ([]*Module, error) {
	files, err := findModuleYmlFiles(root)
	if err != nil {
		return nil, err
	}

	var modules []*Module
	for _, file := range files {
		var module *Module
		module, err = readAndParseModuleYml(file)

		if err != nil {
			return nil, err
		}

		fmt.Printf("Module parsed: %v\n", module.Name)

		modules = append(modules, module)
	}

	return modules, nil
}

func findModuleYmlFiles(root string) ([]string, error) {
	var files []string
	err := fs.WalkDir(fileSystem, root, func(path string, file os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !file.IsDir() && file.Name() == moduleFileName {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func readAndParseModuleYml(filePath string) (*Module, error) {
	data, err := fs.ReadFile(fileSystem, filePath)
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
