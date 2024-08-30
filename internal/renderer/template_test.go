package renderer

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for the text/template interface
type MockTemplate struct {
	mock.Mock
}

func (m *MockTemplate) Execute(wr io.Writer, data interface{}) error {
	args := m.Called(wr, data)

	return args.Error(0)
}

// Tests for the createFileFromTemplate function

func TestCreateFileFromTemplate_MkDirAllReturnsError_ErrorForwarded(t *testing.T) {
	outputFilePath := "test"

	stubMkdirAll := func(path string, perm os.FileMode) error {
		assert.Equal(t, outputFilePath, path)
		assert.Equal(t, perm, os.ModePerm)

		return errors.New("failed to create directory")
	}

	err := createFileFromTemplate(nil, outputFilePath, "", nil, stubMkdirAll, nil)

	assert.Error(t, err)
	assert.Equal(t, "failed to create directory", err.Error())
}

func TestCreateFileFromTemplate_CreateFileReturnsError_ErrorForwarded(t *testing.T) {
	fPath := "test"
	fName := "file.name"

	stubCreateFile := func(path string) (*os.File, error) {
		assert.Equal(t, filepath.Join(fPath, fName), path)

		return nil, errors.New("failed to create file")
	}

	err := createFileFromTemplate(nil, fPath, fName, nil, stubMkdirAll(t, fPath, os.ModePerm), stubCreateFile)

	assert.Error(t, err)
	assert.Equal(t, "failed to create file", err.Error())
}

func TestCreateFileFromTemplate_ExecuteReturnError_ErrorForwarded(t *testing.T) {
	fPath := "test"
	fName := "file.name"

	expectedContext := &context{}
	expectedFile := &os.File{}

	mockTemplate := new(MockTemplate)
	mockTemplate.On("Execute", expectedFile, expectedContext).Return(errors.New("failed to execute template"))

	err := createFileFromTemplate(
		mockTemplate,
		fPath,
		fName,
		expectedContext,
		stubMkdirAll(t, fPath, os.ModePerm),
		stubCreateFile(t, filepath.Join(fPath, fName), expectedFile),
	)

	assert.Error(t, err)
	assert.Equal(t, "failed to execute template", err.Error())
	mockTemplate.AssertExpectations(t)
}

func TestCreateFileFromTemplate_NoErrorsReturned_NilReturned(t *testing.T) {
	fPath := "test"
	fName := "file.name"

	expectedContext := &context{}
	expectedFile := &os.File{}

	mockTemplate := new(MockTemplate)
	mockTemplate.On("Execute", expectedFile, expectedContext).Return(nil)

	err := createFileFromTemplate(
		mockTemplate,
		fPath,
		fName,
		expectedContext,
		stubMkdirAll(t, fPath, os.ModePerm),
		stubCreateFile(t, filepath.Join(fPath, fName), expectedFile),
	)

	assert.Nil(t, err)
	mockTemplate.AssertExpectations(t)
}

// Helper methods

func stubMkdirAll(t *testing.T, expectedPath string, expectedPerm os.FileMode) func(string, os.FileMode) error {
	return func(path string, perm os.FileMode) error {
		assert.Equal(t, expectedPath, path)
		assert.Equal(t, expectedPerm, perm)

		return nil
	}
}

func stubCreateFile(t *testing.T, expectedPath string, returnFile *os.File) func(string) (*os.File, error) {
	return func(path string) (*os.File, error) {
		assert.Equal(t, expectedPath, path)

		return returnFile, nil
	}
}
