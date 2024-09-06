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
	defer setStubMapping(stubMapping)()

	merged, err := MergeWithDefault(stubLanguageName, nil)

	assert.Nil(t, err)
	assert.Equal(t, stubLanguageDefault, merged)
}

func TestMergeWithDefault_PartialLanguageGiven_ReturnsMergedLanguage(t *testing.T) {
	defer setStubMapping(stubMapping)()

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

func TestMergeWithDefault_LanguageWithPartialAdditionalParams_AdditionalParamsMergedWell(t *testing.T) {
	var stubMapping = map[Name]*Language{
		stubLanguageName: {
			AdditionalParameters: map[string]string{
				"myKey":  "defaultValue",
				"myKey2": "defaultValue2",
			},
		},
	}
	defer setStubMapping(stubMapping)()

	actualLanguage := &Language{
		AdditionalParameters: map[string]string{
			"myKey": "actualValue",
		},
	}

	merged, err := MergeWithDefault(stubLanguageName, actualLanguage)

	assert.Nil(t, err)
	assert.Equal(t, map[string]string{
		"myKey":  "actualValue",
		"myKey2": "defaultValue2",
	}, merged.AdditionalParameters)

	// assert default have not been changed
	assert.Equal(t, map[string]string{
		"myKey":  "defaultValue",
		"myKey2": "defaultValue2",
	}, stubMapping[stubLanguageName].AdditionalParameters)

}

func setStubMapping(stubMapping map[Name]*Language) func() {
	originalMapping = defaultMapping
	defaultMapping = stubMapping
	return func() { defaultMapping = originalMapping }
}
