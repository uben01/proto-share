package module

import (
	"io/fs"
	"os"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

var fakeProtoFileContent = []string{
	"syntax = \"proto3\";",
	"package my_package;",
	"message Message {",
	"}",
}

var fakeProtoFileHash = "b0d581adfe0f6f80d63e6169bb35c764"

func TestUpdateModuleVersions_TwoModulesOneWithHashMismatch_HashUpdated(t *testing.T) {
	modules := []*Module{
		{Hash: "123", Version: 1, Path: "module1"},
		{Hash: "456", Version: 1, Path: "module2"},
	}

	defer setComputeModuleMD5Hash(func(fs.FS, string) (string, error) {
		return "123", nil
	})()

	err := UpdateModuleVersions(modules, "")
	assert.Nil(t, err)

	assert.Equal(t, "123", modules[0].Hash)
	assert.Equal(t, 1, modules[0].Version)
	assert.Equal(t, "123", modules[1].Hash)
	assert.Equal(t, 2, modules[1].Version)
}

func TestWriteNewVersionToFile_TwoModulesFirstMarshalReturnsError_ErrorForwarded(t *testing.T) {
	modules := []*Module{
		{Hash: "123", Version: 1, Path: "module1", changed: true},
		{Hash: "321", Version: 1, Path: "module2", changed: true},
	}

	defer setMarshalFunc(func(in interface{}) (out []byte, err error) {
		assert.Equal(t, modules[0], in)

		return nil, assert.AnError
	})()

	err := WriteNewVersionToFile(modules, "")

	assert.Equal(t, assert.AnError, err)
}

func TestWriteNewVersionToFile_ModuleChanged_NewVersionWrittenToFile(t *testing.T) {
	modules := []*Module{
		{Hash: "123", Version: 1, Path: "modules/module1", changed: true},
	}

	defer setMarshalFunc(func(in interface{}) (out []byte, err error) {
		return []byte("marshaled"), nil
	})()
	defer setWriteFileFunc(func(path string, data []byte, perm os.FileMode) error {
		assert.Equal(t, "in/modules/module1/module.yml", path)
		assert.Equal(t, []byte("marshaled"), data)
		assert.Equal(t, os.FileMode(0666), perm)

		return nil
	})()

	err := WriteNewVersionToFile(modules, "in")

	assert.Nil(t, err)
}

func TestWriteNewVersionToFile_ModuleHasNotChanged_DoNothing(t *testing.T) {
	modules := []*Module{
		{Hash: "123", Version: 1, Path: "modules/module1", changed: false},
	}

	defer setMarshalFunc(func(in interface{}) (out []byte, err error) {
		assert.Fail(t, "marshal should not be called")
		return nil, nil
	})()

	err := WriteNewVersionToFile(modules, "in")

	assert.Nil(t, err)
}

func TestComputeFileMD5Hash_FilePresent_FileHashReturned(t *testing.T) {
	filePath := "test.proto"

	testFs := fstest.MapFS{
		filePath: &fstest.MapFile{Data: []byte(strings.Join(fakeProtoFileContent, "\n"))},
	}

	hash, err := computeFileMD5Hash(testFs, filePath)
	assert.Nil(t, err)
	assert.Equal(t, fakeProtoFileHash, hash)
}

func TestComputeFileMD5Hash_FileNotPresent_ErrorReturned(t *testing.T) {
	testFs := fstest.MapFS{}

	hash, err := computeFileMD5Hash(testFs, "test.proto")

	assert.NotNil(t, err)
	assert.Equal(t, "open test.proto: file does not exist", err.Error())
	assert.Empty(t, hash)
}

func TestComputeModuleMD5Hash_NoProtoFilesInDir_ErrorReturned(t *testing.T) {
	testFs := fstest.MapFS{
		"module/test.go": &fstest.MapFile{},
	}

	hash, err := computeModuleMD5Hash(testFs, "module")

	assert.NotNil(t, err)
	assert.Equal(t, "no proto files found in module", err.Error())
	assert.Empty(t, hash)
}

func TestComputeModuleMD5Hash_NoModuleDirFound_ErrorReturned(t *testing.T) {
	testFs := fstest.MapFS{
		"wrongModule/test.go": &fstest.MapFile{},
	}

	hash, err := computeModuleMD5Hash(testFs, "module")

	assert.NotNil(t, err)
	assert.Equal(t, "open module: file does not exist", err.Error())
	assert.Empty(t, hash)
}

func TestComputeModuleMD5Hash_SingleTestFilePresent_HashRehashed(t *testing.T) {
	testFs := fstest.MapFS{
		"module/test.proto": &fstest.MapFile{Data: []byte(strings.Join(fakeProtoFileContent, "\n"))},
	}

	hash, err := computeModuleMD5Hash(testFs, "module/test.proto")

	assert.Nil(t, err)
	assert.Equal(t, "228cbb9f695a2e70fb45efc5888c2021", hash)
}

func setWriteFileFunc(f func(string, []byte, os.FileMode) error) func() {
	originalWriteFile := writeFile
	writeFile = f

	return func() {
		writeFile = originalWriteFile
	}
}

func setMarshalFunc(f func(in interface{}) (out []byte, err error)) func() {
	originalMarshal := marshal
	marshal = f

	return func() {
		marshal = originalMarshal
	}
}

func setComputeModuleMD5Hash(f func(fs.FS, string) (string, error)) func() {
	originalComputeModuleMD5Hash := computeModuleMD5Hash
	computeModuleMD5Hash = f

	return func() {
		computeModuleMD5Hash = originalComputeModuleMD5Hash
	}
}
