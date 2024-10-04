package module

import (
	"io/fs"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/goccy/go-yaml"
)

type Module struct {
	Name    string `yaml:"name"`
	Hash    string `yaml:"hash"`
	Version int    `yaml:"version"`
	Path    string `yaml:"path"`

	Changed bool `yaml:"-"`
}

var moduleFileName = "module.yml"

var fileSystem = os.DirFS(".")

func DiscoverModules(root string) []*Module {
	files := findModuleYmlFiles(root)

	var modules []*Module
	for _, file := range files {
		var module *Module
		module = readAndParseModuleYml(file)

		log.Debugf("Module parsed: %v", module.Name)

		modules = append(modules, module)
	}

	return modules
}

func findModuleYmlFiles(root string) []string {
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

	if err != nil {
		log.Panicf("Error walking the path %q: %v", root, err)
	}

	return files
}

func readAndParseModuleYml(filePath string) *Module {
	data, err := fs.ReadFile(fileSystem, filePath)
	if err != nil {
		log.Panicf("Error reading file %q: %v", filePath, err)
	}

	var module Module
	err = yaml.Unmarshal(data, &module)
	if err != nil {
		log.Panicf("Error parsing file %q: %v", filePath, err)
	}

	return &module
}
