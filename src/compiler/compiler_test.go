package compiler

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"

	. "config"
	. "config/language"
	. "config/module"
)

var (
	testModule = Module{
		Name: "my_module",
		Path: "my/module",
	}

	testLanguage = Language{
		SubDir:             "myLang",
		ModulePathTemplate: "{module}/src/main",
		ProtocCommand:      "myLangOut",
	}

	testConfig = Config{
		InDir:  "indir",
		OutDir: "outdir",
		Modules: []*Module{
			&testModule,
		},
		Languages: map[Name]*Language{
			"myLang": &testLanguage,
		},
	}
)

func TestCompileModules_withoutModules_returnsError(t *testing.T) {
	config := &Config{
		Modules: []*Module{},
	}

	err := CompileModules(config)

	assert.NotNil(t, err)
	assert.Equal(t, "no modules defined", err.Error())
}

func TestCompileModules_withoutLanguages_returnsError(t *testing.T) {
	config := &Config{
		Modules: []*Module{
			{},
		},
		Languages: map[Name]*Language{},
	}

	err := CompileModules(config)

	assert.NotNil(t, err)
	assert.Equal(t, "no languages defined", err.Error())
}

func TestCompileProtos_executeReturnsNil_ErrorReturned(t *testing.T) {
	firstLangOut := "myOut"
	secondLangOut := "myOtherOut"

	expectedProtocCommand := fmt.Sprintf(
		"protoc %s %s -I %s %s/%s/*.proto",
		firstLangOut,
		secondLangOut,
		testConfig.InDir,
		testConfig.InDir,
		testModule.Path,
	)

	stubExecute := func(name string, arg ...string) *exec.Cmd {
		assert.Equal(t, "sh", name)
		assert.Equal(t, []string{"-c", expectedProtocCommand}, arg)

		return nil
	}

	err := compileProtos(&testConfig, &testModule, []string{firstLangOut, secondLangOut}, stubExecute, func(cmd *exec.Cmd) ([]byte, error) { return nil, nil })

	assert.Error(t, err)
	assert.Equal(t, "failed to create command: "+expectedProtocCommand, err.Error())
}

func TestCompileProtos_CombinedOutputReturnsError_ErrorForwarded(t *testing.T) {
	expectedErrorMsg := "stubbed error"

	cmdMock := &exec.Cmd{}
	stubExecute := func(string, ...string) *exec.Cmd {
		return cmdMock
	}

	stubCombinedOutput := func(cmd *exec.Cmd) ([]byte, error) {
		assert.Equal(t, cmdMock, cmd)
		return nil, errors.New(expectedErrorMsg)
	}

	err := compileProtos(&testConfig, &testModule, []string{"myLangOut"}, stubExecute, stubCombinedOutput)

	assert.Error(t, err)
	assert.Equal(t, expectedErrorMsg, err.Error())
}

func TestCompileProtos_CombinedOutputReturnsNoError_NilReturned(t *testing.T) {
	cmdMock := &exec.Cmd{}
	stubExecute := func(string, ...string) *exec.Cmd {
		return cmdMock
	}

	stubCombinedOutput := func(cmd *exec.Cmd) ([]byte, error) {
		return nil, nil
	}

	err := compileProtos(&testConfig, &testModule, []string{"myLangOut"}, stubExecute, stubCombinedOutput)

	assert.Nil(t, err)
}

func TestPrepareLanguageOutput_MkdirAllReturnsError_ErrorForwarded(t *testing.T) {
	expectedErrorMsg := "stubbed error"

	stubMkdirAll := func(string, os.FileMode) error {
		return errors.New(expectedErrorMsg)
	}

	_, err := prepareLanguageOutput(&testConfig, &testLanguage, &testModule, stubMkdirAll)

	assert.Error(t, err)
	assert.Equal(t, expectedErrorMsg, err.Error())
}

func TestPrepareLanguageOutput_LanguagePathTemplateContainsModule_ModuleNameReplaced(t *testing.T) {
	expectedOutputPath := "outdir/myLang/my_module/src/main"

	stubMkdirAll := func(string, os.FileMode) error {
		return nil
	}

	outputPath, err := prepareLanguageOutput(&testConfig, &testLanguage, &testModule, stubMkdirAll)

	assert.Nil(t, err)
	assert.Equal(t, expectedOutputPath, outputPath)
}
