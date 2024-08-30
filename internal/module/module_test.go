package module

import (
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestDiscoverModules_ModuleYamlPresent_ModuleFound(t *testing.T) {
	modulePath := "a/b/module.yml"
	moduleYamlContent := []string{
		"name: test",
		"hash: 123",
		"version: 1",
		"path: a/b",
	}

	testFs := fstest.MapFS{
		modulePath:         &fstest.MapFile{Data: []byte(strings.Join(moduleYamlContent, "\n"))},
		"a/b/module.proto": &fstest.MapFile{},
		"a/module.go":      &fstest.MapFile{},
	}
	defer setFileSystem(testFs)()

	modules, err := DiscoverModules(".")
	assert.Nil(t, err)

	assert.Equal(t, "test", modules[0].Name)
	assert.Equal(t, "123", modules[0].Hash)
	assert.Equal(t, 1, modules[0].Version)
	assert.Equal(t, "a/b", modules[0].Path)
}

func TestDiscoverModules_TwoModulesPresent_ModuleCountMatches(t *testing.T) {
	testFs := fstest.MapFS{
		"a/module.yml": &fstest.MapFile{},
		"b/module.yml": &fstest.MapFile{},
	}
	defer setFileSystem(testFs)()

	modules, err := DiscoverModules(".")
	assert.Nil(t, err)

	assert.Equal(t, 2, len(modules))
}

func setFileSystem(fs fs.FS) func() {
	originalFileSystem := fileSystem
	fileSystem = fs

	return func() {
		fileSystem = originalFileSystem
	}
}
