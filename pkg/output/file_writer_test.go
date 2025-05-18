package output

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewFileOutputWriter_ValidFilePath_001 tests the creation of a new fileWriter instance with a valid file path.
func TestNewFileOutputWriter_ValidFilePath_001(t *testing.T) {
	// Arrange
	tempFile, err := os.CreateTemp("", "testfile")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Act
	writer, err := newFileOutputWriter(tempFile.Name())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, writer)
	assert.NotNil(t, writer.file)
	assert.NotNil(t, writer.writer)

	// Cleanup
	writer.Close()
}

// TestNewFileOutputWriter_InvalidFilePath_002 tests the creation of a new fileWriter instance with an invalid file path.

func TestFileWriter_Write_ValidData_003(t *testing.T) {
	// Arrange
	tempFile, err := os.CreateTemp("", "testfile")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	writer, err := newFileOutputWriter(tempFile.Name())
	require.NoError(t, err)
	defer writer.Close()

	data := []byte("test data")

	// Act
	err = writer.Write(data)

	// Assert
	require.NoError(t, err)

	// Verify written data
	writer.writer.Flush() // Ensure data is flushed to the file
	fileContent, err := os.ReadFile(tempFile.Name())
	require.NoError(t, err)
	assert.Equal(t, string(data)+"\n", string(fileContent))
}

func TestNewFileOutputWriter_InvalidFilePath_002(t *testing.T) {
	// Arrange
	invalidFilePath := "/invalid/path/to/file"

	// Act
	writer, err := newFileOutputWriter(invalidFilePath)

	// Assert
	require.Error(t, err)
	assert.Nil(t, writer)
}

// TestFileWriter_Write_ValidData_003 tests the Write method of fileWriter with valid data.

