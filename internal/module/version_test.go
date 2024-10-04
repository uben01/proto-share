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

	defer setComputeModuleMD5Hash(func(fs.FS, string) string {
		return "123"
	})()

	UpdateModuleVersions(modules, "")

	assert.Equal(t, "123", modules[0].Hash)
	assert.Equal(t, 1, modules[0].Version)
	assert.Equal(t, "123", modules[1].Hash)
	assert.Equal(t, 2, modules[1].Version)
}

func TestWriteNewVersionToFile_TwoModulesFirstMarshalReturnsError_Panics(t *testing.T) {
	modules := []*Module{
		{Hash: "123", Version: 1, Path: "module1", Changed: true},
		{Hash: "321", Version: 1, Path: "module2", Changed: true},
	}

	defer setMarshalFunc(func(in interface{}) (out []byte, err error) {
		assert.Equal(t, modules[0], in)

		return nil, assert.AnError
	})()

	assert.Panics(t, func() {
		WriteNewVersionToFile(modules, "")
	})
}

func TestWriteNewVersionToFile_ModuleChanged_NewVersionWrittenToFile(t *testing.T) {
	modules := []*Module{
		{Hash: "123", Version: 1, Path: "modules/module1", Changed: true},
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

	WriteNewVersionToFile(modules, "in")
}

func TestWriteNewVersionToFile_ModuleHasNotChanged_DoNothing(t *testing.T) {
	modules := []*Module{
		{Hash: "123", Version: 1, Path: "modules/module1", Changed: false},
	}

	defer setMarshalFunc(func(in interface{}) (out []byte, err error) {
		assert.Fail(t, "marshal should not be called")
		return nil, nil
	})()

	WriteNewVersionToFile(modules, "in")
}

func TestComputeFileMD5Hash_FilePresent_FileHashReturned(t *testing.T) {
	filePath := "test.proto"

	testFs := fstest.MapFS{
		filePath: &fstest.MapFile{Data: []byte(strings.Join(fakeProtoFileContent, "\n"))},
	}

	hash := computeFileMD5Hash(testFs, filePath)
	assert.Equal(t, fakeProtoFileHash, hash)
}

func TestComputeFileMD5Hash_FileNotPresent_Panics(t *testing.T) {
	testFs := fstest.MapFS{}

	assert.Panics(t, func() {
		_ = computeFileMD5Hash(testFs, "test.proto")
	})
}

func TestComputeModuleMD5Hash_NoProtoFilesInDir_Panics(t *testing.T) {
	testFs := fstest.MapFS{
		"module/test.go": &fstest.MapFile{},
	}

	assert.Panics(t, func() {
		_ = computeModuleMD5Hash(testFs, "module")
	})
}

func TestComputeModuleMD5Hash_NoModuleDirFound_Panics(t *testing.T) {
	testFs := fstest.MapFS{
		"wrongModule/test.go": &fstest.MapFile{},
	}

	assert.Panics(t, func() {
		_ = computeModuleMD5Hash(testFs, "module")
	})
}

func TestComputeModuleMD5Hash_SingleTestFilePresent_HashRehashed(t *testing.T) {
	testFs := fstest.MapFS{
		"module/test.proto": &fstest.MapFile{Data: []byte(strings.Join(fakeProtoFileContent, "\n"))},
	}

	hash := computeModuleMD5Hash(testFs, "module/test.proto")

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

func setComputeModuleMD5Hash(f func(fs.FS, string) string) func() {
	originalComputeModuleMD5Hash := computeModuleMD5Hash
	computeModuleMD5Hash = f

	return func() {
		computeModuleMD5Hash = originalComputeModuleMD5Hash
	}
}
