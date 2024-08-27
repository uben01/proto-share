package language

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const mockLanguageName Name = "myLanguage"

func mockLanguageDefault() *Language {
	return &Language{
		ModuleTemplatePath:   "myModulePathTemplate",
		SubDir:               "mySubDir",
		ProtocCommand:        "myProtocCommand",
		AdditionalParameters: map[string]string{"myKey": "myValue"},
	}
}

var originalMapping = defaultMapping
var mockMapping = map[Name]*Language{
	mockLanguageName: mockLanguageDefault(),
}

func TestMergeWithDefault_NoLanguageMapping_ReturnsError(t *testing.T) {
	defaultMapping = map[Name]*Language{}
	defer func() { defaultMapping = originalMapping }()

	_, err := MergeWithDefault(mockLanguageName, nil)

	assert.NotNil(t, err)
	assert.Equal(t, "unsupported language: mylanguage", err.Error())
}

func TestMergeWithDefault_NoActualLanguage_ReturnsDefaultLanguage(t *testing.T) {
	defaultMapping = mockMapping
	defer func() { defaultMapping = originalMapping }()

	merged, err := MergeWithDefault(mockLanguageName, nil)

	assert.Nil(t, err)
	assert.Equal(t, mockLanguageDefault(), merged)
}
