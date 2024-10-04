package config

import (
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	. "github.com/uben01/proto-share/internal/language"
)

type MergeWithDefaultMock struct {
	mock.Mock
}

func (m *MergeWithDefaultMock) MergeWithDefault(languageName Name, actualLanguage *Language) *Language {
	args := m.Called(languageName, actualLanguage)

	return args.Get(0).(*Language)
}

func TestParseConfig_FileNotFound_Panics(t *testing.T) {
	defer setFileSystem(fstest.MapFS{})()

	assert.Panics(t, func() {
		_ = ParseConfig("config.yaml")
	})
}

func TestParseConfig_ConfigContainsLanguages_MergeWithDefaultCalledForEach(t *testing.T) {
	modulePath := "config.yaml"
	moduleYamlContent := []string{
		"projectName: test",
		"inDir: in",
		"outDir: out",
		"languages:",
		"  Java:",
	}

	testFs := fstest.MapFS{
		modulePath: &fstest.MapFile{Data: []byte(strings.Join(moduleYamlContent, "\n"))},
	}
	defer setFileSystem(testFs)()
	mockMergeWithDefault := &MergeWithDefaultMock{}
	defer setMergeWithDefault(mockMergeWithDefault)()

	stubLanguage := &Language{}
	mockMergeWithDefault.
		On("MergeWithDefault", Java, mock.Anything).
		Once().
		Return(stubLanguage, nil)

	config := ParseConfig(modulePath)

	assert.Equal(t, 1, len(config.Languages))
	assert.Equal(t, stubLanguage, config.Languages[Java])
	assert.Equal(t, "test", config.ProjectName)
	assert.Equal(t, "in", config.InDir)
	assert.Equal(t, "out", config.OutDir)

	mockMergeWithDefault.AssertNumberOfCalls(t, "MergeWithDefault", 1)
	mockMergeWithDefault.AssertExpectations(t)
}

func setFileSystem(fs fs.FS) func() {
	originalFileSystem := fileSystem
	fileSystem = fs

	return func() {
		fileSystem = originalFileSystem
	}
}

func setMergeWithDefault(mergeWithDefaultMock *MergeWithDefaultMock) func() {
	originalMergeWithDefault := mergeWithDefault
	mergeWithDefault = mergeWithDefaultMock.MergeWithDefault

	return func() {
		mergeWithDefault = originalMergeWithDefault
	}
}
