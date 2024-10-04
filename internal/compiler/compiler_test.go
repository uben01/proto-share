package compiler

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/uben01/proto-share/internal/config"
	. "github.com/uben01/proto-share/internal/context"
	. "github.com/uben01/proto-share/internal/language"
	. "github.com/uben01/proto-share/internal/module"
)

var (
	testModule = Module{
		Name: "my_module",
		Path: "my/module",
	}

	testLanguage = Language{
		SubDir:            "myLang",
		ModuleCompilePath: "{{ .Module.Name }}/src/main",
		ProtocCommand:     "myLangOut",
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

func TestCompileModules_withoutModules_Panics(t *testing.T) {
	config := &Config{
		Modules: []*Module{},
	}

	assert.Panics(t, func() {
		CompileModules(config)
	})
}

func TestCompileModules_withoutLanguages_Panics(t *testing.T) {
	config := &Config{
		Modules: []*Module{
			{},
		},
		Languages: map[Name]*Language{},
	}

	assert.Panics(t, func() {
		CompileModules(config)
	})
}

func TestCompileModules_changedIsFalse_prepareLanguageOutputNotCalled(t *testing.T) {
	config := &Config{
		Modules: []*Module{
			{Changed: false},
		},
		Languages: map[Name]*Language{
			"": {},
		},
	}

	defer setStubPrepareLanguageOutput(func(*Config, *Language, func(string, os.FileMode) error) string {
		assert.Fail(t, "prepareLanguageOutput should not be called")
		return ""
	})()

	CompileModules(config)
}

func TestCompileProtos_executeReturnsNil_Panics(t *testing.T) {
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

	assert.Panics(t, func() {
		compileProtos(
			&testConfig,
			&testModule,
			[]string{firstLangOut, secondLangOut},
			stubExecute,
			func(cmd *exec.Cmd) ([]byte, error) { return nil, nil },
		)
	})
}

func TestCompileProtos_CombinedOutputReturnsError_Panics(t *testing.T) {
	expectedErrorMsg := "stubbed error"

	cmdMock := &exec.Cmd{}
	stubExecute := func(string, ...string) *exec.Cmd {
		return cmdMock
	}

	stubCombinedOutput := func(cmd *exec.Cmd) ([]byte, error) {
		assert.Equal(t, cmdMock, cmd)
		return nil, errors.New(expectedErrorMsg)
	}

	assert.Panics(t, func() {
		compileProtos(&testConfig, &testModule, []string{"myLangOut"}, stubExecute, stubCombinedOutput)
	})
}

func TestCompileProtos_CombinedOutputReturnsNoError_notPanics(t *testing.T) {
	cmdMock := &exec.Cmd{}
	stubExecute := func(string, ...string) *exec.Cmd {
		return cmdMock
	}

	stubCombinedOutput := func(cmd *exec.Cmd) ([]byte, error) {
		return nil, nil
	}

	compileProtos(&testConfig, &testModule, []string{"myLangOut"}, stubExecute, stubCombinedOutput)
}

func TestPrepareLanguageOutput_MkdirAllReturnsError_Panics(t *testing.T) {
	CTX = &Context{Language: &testLanguage, Module: &testModule}

	expectedErrorMsg := "stubbed error"

	stubMkdirAll := func(string, os.FileMode) error {
		return errors.New(expectedErrorMsg)
	}

	assert.Panics(t, func() {
		prepareLanguageOutput(&testConfig, &testLanguage, stubMkdirAll)
	})
}

func TestPrepareLanguageOutput_LanguagePathTemplateContainsModule_ModuleNameReplaced(t *testing.T) {
	expectedOutputPath := "outdir/myLang/my_module/src/main"

	stubMkdirAll := func(string, os.FileMode) error {
		return nil
	}

	outputPath := prepareLanguageOutput(&testConfig, &testLanguage, stubMkdirAll)

	assert.Equal(t, expectedOutputPath, outputPath)
}

func setStubPrepareLanguageOutput(f func(*Config, *Language, func(string, os.FileMode) error) string) func() {
	originalPrepareLanguageOutput := prepareLanguageOutput
	prepareLanguageOutput = f
	return func() { prepareLanguageOutput = originalPrepareLanguageOutput }
}
