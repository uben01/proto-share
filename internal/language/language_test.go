package language

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const stubLanguageName Name = "myLanguage"

var stubLanguageDefault = &Language{
	ModuleTemplatePath:   "myModulePathTemplate",
	SubDir:               "mySubDir",
	ProtocCommand:        "myProtocCommand",
	AdditionalParameters: map[string]string{"myKey": "myValue"},
}

var originalMapping = defaultMapping
var stubMapping = map[Name]*Language{
	stubLanguageName: stubLanguageDefault,
}

func TestMergeWithDefault_NoLanguageMapping_ReturnsError(t *testing.T) {
	_, err := MergeWithDefault(stubLanguageName, nil)

	assert.NotNil(t, err)
	assert.Equal(t, "unsupported language: mylanguage", err.Error())
}

func TestMergeWithDefault_NoActualLanguage_ReturnsDefaultLanguage(t *testing.T) {
	defer setStubMapping()()

	merged, err := MergeWithDefault(stubLanguageName, nil)

	assert.Nil(t, err)
	assert.Equal(t, stubLanguageDefault, merged)
}

func TestMergeWithDefault_PartialLanguageGiven_ReturnsMergedLanguage(t *testing.T) {
	defer setStubMapping()()

	actualLanguage := &Language{
		SubDir: "actualSubDir",
	}

	merged, err := MergeWithDefault(stubLanguageName, actualLanguage)

	assert.Nil(t, err)
	assert.Equal(t, "actualSubDir", merged.SubDir)
	assert.Equal(t, "myModulePathTemplate", merged.ModuleTemplatePath)
	assert.Equal(t, "myProtocCommand", merged.ProtocCommand)
	assert.Equal(t, map[string]string{"myKey": "myValue"}, merged.AdditionalParameters)
}

func setStubMapping() func() {
	originalMapping = defaultMapping
	defaultMapping = stubMapping
	return func() { defaultMapping = originalMapping }
}
