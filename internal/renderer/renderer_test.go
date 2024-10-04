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

func (m *MockWalkTemplateDir) walkTemplateDir(fs fs.FS, from string, to string) {
	m.Called(fs, from, to)
}

// Tests for RenderTemplates function

func TestRenderTemplates_NoLanguages_Panics(t *testing.T) {
	conf := &config.Config{}

	testFs := fstest.MapFS{}

	assert.Panics(t, func() {
		RenderTemplates(testFs, conf)
	})
}

func TestRenderTemplates_NoModules_Panics(t *testing.T) {
	conf := &config.Config{
		Languages: map[language.Name]*language.Language{
			"test": {},
		},
	}

	testFs := fstest.MapFS{}

	assert.Panics(t, func() {
		RenderTemplates(testFs, conf)
	})
}

func TestRenderTemplates_multipleLanguagesAndModules_walkTemplateDirCalledForEveryCombination(t *testing.T) {
	conf := &config.Config{
		OutDir: "out",
		Languages: map[language.Name]*language.Language{
			"languageName1": {SubDir: "languagename1", ModuleTemplatePath: "module"},
			"languageName2": {SubDir: "languagename2", ModuleTemplatePath: "{{ .Module.Name }}"},
		}, Modules: []*module.Module{
			{Name: "module1", Changed: true},
			{Name: "module2", Changed: true},
			{Name: "module2", Changed: false},
		},
	}

	CTX = &Context{Config: conf}

	mockWalkTemplateDir := new(MockWalkTemplateDir)
	setWalkTemplateDirFunc(mockWalkTemplateDir)

	testFs := fstest.MapFS{}
	// Language globals
	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename1/global", "out/languagename1").
		Once()
	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename2/global", "out/languagename2").
		Once()

	// Module for languages
	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename1/module", "out/languagename1/module").
		Once()

	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename1/module", "out/languagename1/module").
		Once()

	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename2/module", "out/languagename2/module1").
		Once()

	mockWalkTemplateDir.
		On("walkTemplateDir", testFs, "assets/templates/languagename2/module", "out/languagename2/module2").
		Once()

	RenderTemplates(testFs, conf)

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
