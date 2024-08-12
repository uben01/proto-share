package compiler

import (
	"errors"
	"fmt"
	"os"
	. "os/exec"
	"testing"

	"github.com/stretchr/testify/assert"

	. "config"
	"config/language"
	. "config/module"
)

var (
	testModule = Module{
		Name: "my_module",
		Path: "my/module",
	}

	testLanguage = language.Config{
		SubDir:         "myLang",
		ModulePath:     "myModuleGoesHere",
		ProtoOutputDir: "src/main",
		ProtocCommand:  "myLangOut",
	}

	testConfig = Config{
		InDir:  "indir",
		OutDir: "outdir",
		Modules: []*Module{
			&testModule,
		},
		Languages: map[language.Name]*language.Config{
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
		Languages: map[language.Name]*language.Config{},
	}

	err := CompileModules(config)

	assert.NotNil(t, err)
	assert.Equal(t, "no languages defined", err.Error())
}

func TestCompileModules_createLayoutReturnsError_errorForwarded(t *testing.T) {
	defer stubMkDirAll(func(path string, perm os.FileMode) error {
		expectedPath := fmt.Sprintf(
			"%s/%s/%s/%s/%s",
			testConfig.OutDir,
			testLanguage.SubDir,
			testLanguage.ModulePath,
			testModule.Name,
			testLanguage.ProtoOutputDir,
		)
		assert.Equal(t, expectedPath, path)
		assert.Equal(t, os.ModePerm, perm)

		return errors.New("stubbed error")
	})()

	err := CompileModules(&testConfig)

	assert.NotNil(t, err)
	assert.Equal(t, "stubbed error", err.Error())
}

func TestCompileModules_createCommandReturnsNil_errorReturned(t *testing.T) {
	defer stubMkDirAll(func(string, os.FileMode) error {
		return nil
	})()

	expectedProtocCommand := fmt.Sprintf(
		"protoc --%s=%s/%s/%s/%s/%s -I %s %s/%s/*.proto",
		testLanguage.ProtocCommand,
		testConfig.OutDir,
		testLanguage.SubDir,
		testLanguage.ModulePath,
		testModule.Name,
		testLanguage.ProtoOutputDir,
		testConfig.InDir,
		testConfig.InDir,
		testModule.Path,
	)

	defer stubCreateCommand(func(name string, arg ...string) *Cmd {
		assert.Equal(t, "sh", name)
		assert.Len(t, arg, 2)
		assert.Equal(t, "-c", arg[0])

		assert.Equal(t, expectedProtocCommand, arg[1])

		return nil
	})()

	err := CompileModules(&testConfig)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("failed to create command: %s", expectedProtocCommand), err.Error())
}

func TestCompileModules_executeCommandReturnsError_errorReturned(t *testing.T) {
	defer stubMkDirAll(func(string, os.FileMode) error {
		return nil
	})()

	stubCmd := &Cmd{}
	defer stubCreateCommand(func(string, ...string) *Cmd {
		return stubCmd
	})()

	defer stubExecuteCommand(func(cmd *Cmd) ([]byte, error) {
		assert.Equal(t, stubCmd, cmd)

		return nil, errors.New("stubbed error")
	})()

	err := CompileModules(&testConfig)

	assert.NotNil(t, err)
	assert.Equal(t, "stubbed error", err.Error())
}

func TestCompileModules_noError_returnsNil(t *testing.T) {
	defer stubMkDirAll(func(string, os.FileMode) error {
		return nil
	})()

	defer stubCreateCommand(func(string, ...string) *Cmd {
		return &Cmd{}
	})()

	defer stubExecuteCommand(func(*Cmd) ([]byte, error) {
		return nil, nil
	})()

	err := CompileModules(&testConfig)

	assert.Nil(t, err)
}

func stubMkDirAll(stub func(string, os.FileMode) error) func() {
	originalMkDirAll := mkdirAll

	mkdirAll = stub

	return func() {
		mkdirAll = originalMkDirAll
	}
}

func stubCreateCommand(stub func(string, ...string) *Cmd) func() {
	originalCreateCommand := createCommand

	createCommand = stub

	return func() {
		createCommand = originalCreateCommand
	}
}

func stubExecuteCommand(stub func(*Cmd) ([]byte, error)) func() {
	originalExecuteCommand := executeCommand

	executeCommand = stub

	return func() {
		executeCommand = originalExecuteCommand
	}
}
