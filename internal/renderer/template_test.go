package renderer

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
	templ "text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/uben01/proto-share/internal/config"
	"github.com/uben01/proto-share/internal/language"
	"github.com/uben01/proto-share/internal/module"
)

// Mock for walkTemplateDir function
type MockWalkTemplateDir struct {
	mock.Mock
}

func (m *MockWalkTemplateDir) walkTemplateDir(fs fs.FS, from string, to string, context *context) error {
	args := m.Called(fs, from, to, context)

	return args.Error(0)
}

// Tests for the createFileFromTemplate function

func TestCreateFileFromTemplate_MkDirAllReturnsError_ErrorForwarded(t *testing.T) {
	outputFilePath := "test"

	stubMkdirAll := func(path string, perm os.FileMode) error {
		assert.Equal(t, outputFilePath, path)
		assert.Equal(t, perm, os.ModePerm)

		return errors.New("failed to create directory")
	}

	err := createFileFromTemplate(nil, outputFilePath, "", nil, stubMkdirAll, nil)

	assert.Error(t, err)
	assert.Equal(t, "failed to create directory", err.Error())
}

func TestCreateFileFromTemplate_CreateFileReturnsError_ErrorForwarded(t *testing.T) {
	fPath := "test"
	fName := "file.name"

	stubCreateFile := func(path string) (*os.File, error) {
		assert.Equal(t, filepath.Join(fPath, fName), path)

		return nil, errors.New("failed to create file")
	}

	err := createFileFromTemplate(nil, fPath, fName, nil, stubMkdirAll(t, fPath, os.ModePerm), stubCreateFile)

	assert.Error(t, err)
	assert.Equal(t, "failed to create file", err.Error())
}

// Tests for RenderTemplates function

func TestRenderTemplates_NoLanguages_ReturnsError(t *testing.T) {
	conf := &config.Config{}

	testFs := fstest.MapFS{}

	err := RenderTemplates(testFs, conf)

	assert.Error(t, err)
	assert.Equal(t, "no languages found in config", err.Error())
}

func TestRenderTemplates_NoModules_ReturnsError(t *testing.T) {
	conf := &config.Config{
		Languages: map[language.Name]*language.Language{
			"test": {},
		},
	}

	testFs := fstest.MapFS{}

	err := RenderTemplates(testFs, conf)

	assert.Error(t, err)
	assert.Equal(t, "no modules found in config", err.Error())
}

func TestRenderTemplates_multipleLanguagesAndModules_walkTemplateDirCalledForEveryCombination(t *testing.T) {
	conf := &config.Config{
		OutDir: "out",
		Languages: map[language.Name]*language.Language{
			"languageName1": {SubDir: "languagename1", ModuleTemplatePath: "module"},
			"languageName2": {SubDir: "languagename2", ModuleTemplatePath: "{module}"},
		}, Modules: []*module.Module{
			{Name: "module1"},
			{Name: "module2"},
		},
	}

	ctx := &context{Config: conf}
	defer setNewContextFunc(ctx)()

	mockWalkTemplateDir := new(MockWalkTemplateDir)
	setWalkTemplateDirFunc(mockWalkTemplateDir)

	testFs := fstest.MapFS{}
	// Language globals
	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename1/global", "out/languagename1", ctx).
		Once().
		Return(nil)
	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename2/global", "out/languagename2", ctx).
		Once().
		Return(nil)

	// Module for languages
	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename1/module", "out/languagename1/module", ctx).
		Once().
		Return(nil)

	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename1/module", "out/languagename1/module", ctx).
		Once().
		Return(nil)

	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename2/module", "out/languagename2/module1", ctx).
		Once().
		Return(nil)

	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename2/module", "out/languagename2/module2", ctx).
		Once().
		Return(nil)

	err := RenderTemplates(testFs, conf)

	assert.Nil(t, err)

	mockWalkTemplateDir.AssertNumberOfCalls(t, "walkTemplateDir", 6)
	mockWalkTemplateDir.AssertExpectations(t)
}

// Tests for processTemplateRecursively function

func TestProcessTemplateRecursively_OneLevel_OutputReturned(t *testing.T) {
	templateString := `{{ .Module.Path }}`
	context := &context{
		Module: &module.Module{
			Path: "module",
		},
	}

	template, _ := templ.New("").Funcs(customFunctions).Parse(templateString)

	output, err := processTemplateRecursively(template, context)
	assert.Nil(t, err)
	assert.Equal(t, "module", output)
}

func TestProcessTemplateRecursively_TwoLevelsWithCustomFunctions_OutputReturned(t *testing.T) {
	templateString := `{{ .Module.Path }}`
	context := &context{
		Module: &module.Module{
			Name: "moduleName",
			Path: "{{ .Module.Name | kebabCase }}",
		},
	}

	template, _ := templ.New("").Funcs(customFunctions).Parse(templateString)

	output, err := processTemplateRecursively(template, context)
	assert.Nil(t, err)
	assert.Equal(t, "module-name", output)
}

// Helper methods

func stubMkdirAll(t *testing.T, expectedPath string, expectedPerm os.FileMode) func(string, os.FileMode) error {
	return func(path string, perm os.FileMode) error {
		assert.Equal(t, expectedPath, path)
		assert.Equal(t, expectedPerm, perm)

		return nil
	}
}

func setNewContextFunc(c *context) func() {
	originalNewContext := newContext
	newContext = func(*config.Config) *context {
		return c
	}

	return func() {
		newContext = originalNewContext
	}
}

func setWalkTemplateDirFunc(m *MockWalkTemplateDir) func() {
	originalWalkTemplateDir := walkTemplateDir
	walkTemplateDir = m.walkTemplateDir

	return func() {
		walkTemplateDir = originalWalkTemplateDir
	}
}
