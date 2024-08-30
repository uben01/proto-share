package renderer

import (
	"errors"
	"io"
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

// Mock for the text/template interface
type MockTemplate struct {
	mock.Mock
}

func (m *MockTemplate) Execute(wr io.Writer, data interface{}) error {
	args := m.Called(wr, data)

	return args.Error(0)
}

// Mock for createFileFromTemplate function
type MockCreateFileFromTemplate struct {
	mock.Mock
}

func (m *MockCreateFileFromTemplate) createFileFromTemplate(
	tmpl templateExecutor,
	outputPath string,
	outputFileName string,
	context *context,
	mkdirAll func(string, os.FileMode) error,
	createFile func(string) (*os.File, error),
) error {
	args := m.Called(tmpl, outputPath, outputFileName, context, mkdirAll, createFile)

	return args.Error(0)
}

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

func TestCreateFileFromTemplate_ExecuteReturnError_ErrorForwarded(t *testing.T) {
	fPath := "test"
	fName := "file.name"

	expectedContext := &context{}
	expectedFile := &os.File{}

	mockTemplate := new(MockTemplate)
	mockTemplate.On("Execute", expectedFile, expectedContext).Return(errors.New("failed to execute template"))

	err := createFileFromTemplate(
		mockTemplate,
		fPath,
		fName,
		expectedContext,
		stubMkdirAll(t, fPath, os.ModePerm),
		stubCreateFile(t, filepath.Join(fPath, fName), expectedFile),
	)

	assert.Error(t, err)
	assert.Equal(t, "failed to execute template", err.Error())
	mockTemplate.AssertExpectations(t)
}

func TestCreateFileFromTemplate_NoErrorsReturned_NilReturned(t *testing.T) {
	fPath := "test"
	fName := "file.name"

	expectedContext := &context{}
	expectedFile := &os.File{}

	mockTemplate := new(MockTemplate)
	mockTemplate.On("Execute", expectedFile, expectedContext).Return(nil)

	err := createFileFromTemplate(
		mockTemplate,
		fPath,
		fName,
		expectedContext,
		stubMkdirAll(t, fPath, os.ModePerm),
		stubCreateFile(t, filepath.Join(fPath, fName), expectedFile),
	)

	assert.Nil(t, err)
	mockTemplate.AssertExpectations(t)
}

// Tests for the walkTemplateDir function

func TestWalkTemplateDir(t *testing.T) {
	context := &context{}
	template := &templ.Template{}

	mockCreateFileFromTemplate := new(MockCreateFileFromTemplate)
	defer setCreateFileFromTemplateFunc(mockCreateFileFromTemplate)()
	defer setParseTemplateFunc(template, nil)()

	mockCreateFileFromTemplate.
		On("createFileFromTemplate", template, "to", "file1.tmpl", context, mock.Anything, mock.Anything).
		Once().
		Return(nil)
	mockCreateFileFromTemplate.
		On("createFileFromTemplate", template, "to/dir1", "file2.tmpl", context, mock.Anything, mock.Anything).
		Once().
		Return(nil)
	mockCreateFileFromTemplate.
		On("createFileFromTemplate", template, "to/dir2", "file3.tmpl", context, mock.Anything, mock.Anything).
		Once().
		Return(nil)

	testFs := fstest.MapFS{
		"from/file1.tmpl":      &fstest.MapFile{},
		"from/dir1/file2.tmpl": &fstest.MapFile{},
		"from/dir2/file3.tmpl": &fstest.MapFile{},
	}

	err := walkTemplateDir(testFs, "from", "to", context)

	assert.NoError(t, err)
	mockCreateFileFromTemplate.AssertNumberOfCalls(t, "createFileFromTemplate", 3)
	mockCreateFileFromTemplate.AssertExpectations(t)
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
	conf := &config.Config{}
	conf.Languages = map[language.Name]*language.Language{
		"test": {},
	}

	testFs := fstest.MapFS{}

	err := RenderTemplates(testFs, conf)

	assert.Error(t, err)
	assert.Equal(t, "no modules found in config", err.Error())
}

func TestRenderTemplates_multipleLanguagesAndModules_walkTemplateDirCalledForEveryCombination(t *testing.T) {
	conf := &config.Config{}
	conf.OutDir = "out"
	conf.Languages = map[language.Name]*language.Language{
		"languageName1": {SubDir: "languagename1", ModuleTemplatePath: "module"},
		"languageName2": {SubDir: "languagename2", ModuleTemplatePath: "{module}"},
	}
	conf.Modules = []*module.Module{
		{Name: "module1"},
		{Name: "module2"},
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

// Helper methods

func stubMkdirAll(t *testing.T, expectedPath string, expectedPerm os.FileMode) func(string, os.FileMode) error {
	return func(path string, perm os.FileMode) error {
		assert.Equal(t, expectedPath, path)
		assert.Equal(t, expectedPerm, perm)

		return nil
	}
}

func stubCreateFile(t *testing.T, expectedPath string, returnFile *os.File) func(string) (*os.File, error) {
	return func(path string) (*os.File, error) {
		assert.Equal(t, expectedPath, path)

		return returnFile, nil
	}
}

func setCreateFileFromTemplateFunc(m *MockCreateFileFromTemplate) func() {
	originalCreateFileFromTemplate := createFileFromTemplate
	createFileFromTemplate = m.createFileFromTemplate

	return func() {
		createFileFromTemplate = originalCreateFileFromTemplate
	}
}

func setParseTemplateFunc(template *templ.Template, err error) func() {
	originalParseTemplate := parseTemplate
	parseTemplate = func(fs fs.FS, patterns ...string) (*templ.Template, error) {
		return template, err
	}

	return func() {
		parseTemplate = originalParseTemplate
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
