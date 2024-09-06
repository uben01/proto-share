package renderer

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/uben01/proto-share/internal/config"
	. "github.com/uben01/proto-share/internal/context"
	"github.com/uben01/proto-share/internal/language"
	"github.com/uben01/proto-share/internal/module"
)

// Mock for walkTemplateDir function
type MockWalkTemplateDir struct {
	mock.Mock
}

func (m *MockWalkTemplateDir) walkTemplateDir(fs fs.FS, from string, to string) error {
	args := m.Called(fs, from, to)

	return args.Error(0)
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
			"languageName2": {SubDir: "languagename2", ModuleTemplatePath: "{{ .Module.Name }}"},
		}, Modules: []*module.Module{
			{Name: "module1"},
			{Name: "module2"},
		},
	}

	CTX = &Context{Config: conf}

	mockWalkTemplateDir := new(MockWalkTemplateDir)
	setWalkTemplateDirFunc(mockWalkTemplateDir)

	testFs := fstest.MapFS{}
	// Language globals
	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename1/global", "out/languagename1").
		Once().
		Return(nil)
	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename2/global", "out/languagename2").
		Once().
		Return(nil)

	// Module for languages
	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename1/module", "out/languagename1/module").
		Once().
		Return(nil)

	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename1/module", "out/languagename1/module").
		Once().
		Return(nil)

	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename2/module", "out/languagename2/module1").
		Once().
		Return(nil)

	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename2/module", "out/languagename2/module2").
		Once().
		Return(nil)

	err := RenderTemplates(testFs, conf)

	assert.Nil(t, err)

	mockWalkTemplateDir.AssertNumberOfCalls(t, "walkTemplateDir", 6)
	mockWalkTemplateDir.AssertExpectations(t)
}

// Helper methods

func setWalkTemplateDirFunc(m *MockWalkTemplateDir) func() {
	originalWalkTemplateDir := walkTemplateDir
	walkTemplateDir = m.walkTemplateDir

	return func() {
		walkTemplateDir = originalWalkTemplateDir
	}
}
