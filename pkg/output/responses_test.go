package output

import (
	"crypto/sha1"
	"encoding/hex"
	"testing"

	"os"
	"path/filepath"

	"net/http"

	"github.com/projectdiscovery/katana/pkg/navigation"
	"github.com/stretchr/testify/assert"
)

// TestGetResponseHash_ValidURL_123 tests the `getResponseHash` function with a valid URL input.
func TestUpdateIndex_ValidInput_654(t *testing.T) {
	// Arrange
	storeResponseFolder := "/tmp/responses"
	indexFilePath := filepath.Join(storeResponseFolder, indexFile)
	_ = os.MkdirAll(storeResponseFolder, os.ModePerm)
	_ = os.WriteFile(indexFilePath, []byte{}, 0644)

	result := &Result{
		Request: &navigation.Request{
			URL: "http://example.com/path",
		},
		Response: &navigation.Response{
			Resp: &http.Response{Status: "200 OK"},
		},
	}

	// Act
	err := updateIndex(storeResponseFolder, result)

	// Assert
	assert.NoError(t, err, "Error should be nil for valid inputs")
	content, readErr := os.ReadFile(indexFilePath)
	assert.NoError(t, readErr, "Error reading index file should be nil")
	expectedContent := getResponseFileName(storeResponseFolder, "example.com", result.Request.URL) +
		" " + result.Request.URL + " (200 OK)\n"
	assert.Contains(t, string(content), expectedContent, "Index file content should match the expected value")
}

// TestGetResponseFile_ValidInput_654 tests the `getResponseFile` function with valid inputs.

func TestFormatResult_CorrectOutput_987(t *testing.T) {
	// Arrange
	writer := &StandardWriter{}
	result := &Result{
		Request: &navigation.Request{
			URL: "http://example.com",
			Raw: "GET / HTTP/1.1\nHost: example.com\n\n", // Raw request has trailing newlines
		},
		Response: &navigation.Response{
			Raw: "HTTP/1.1 200 OK\nContent-Type: text/html\n\n<html></html>",
		},
	}

	// Act
	formattedResult, err := writer.formatResult(result)

	// Assert
	assert.NoError(t, err, "Error should be nil for valid inputs")
	// Expected: URL + "\n\n\n" + Request.Raw + "\n\n" + Response.Raw
	// Given Request.Raw ends with "\n\n", the concatenation results in 4 newlines between Raw request and Raw response.
	expectedResultString := "http://example.com\n\n\nGET / HTTP/1.1\nHost: example.com\n\n\n\nHTTP/1.1 200 OK\nContent-Type: text/html\n\n<html></html>"
	assert.Equal(t, expectedResultString, string(formattedResult), "Formatted result should match the expected value")
}

func TestGetResponseFile_ValidInput_654(t *testing.T) {
	// Arrange
	storeResponseFolder := "/tmp/responses"
	url := "http://example.com/path"
	domain := "example.com"
	expectedFileName := filepath.Join(storeResponseFolder, domain, getResponseHash(url)+".txt")
	_ = os.MkdirAll(filepath.Join(storeResponseFolder, domain), os.ModePerm)

	// Act
	fileName, fileWriter, err := getResponseFile(storeResponseFolder, url)

	// Assert
	assert.NoError(t, err, "Error should be nil for valid inputs")
	assert.Equal(t, expectedFileName, fileName, "File name should match the expected value")
	assert.NotNil(t, fileWriter, "File writer should not be nil")
}

// TestFormatResult_CorrectOutput_987 ensures formatResult produces the correct string output.
// It specifically validates the newline characters based on the function's logic.

func TestGetResponseHost_ValidURL_456(t *testing.T) {
	// Arrange
	url := "http://example.com:8080/path"
	expectedHost := "example.com_8080"

	// Act
	actualHost, err := getResponseHost(url)

	// Assert
	assert.NoError(t, err, "Error should be nil for valid URL")
	assert.Equal(t, expectedHost, actualHost, "Host should match the expected value")
}

// TestCreateHostDir_ValidInput_789 tests the `createHostDir` function with valid inputs.

func TestGetResponseFileName_ValidInput_321(t *testing.T) {
	// Arrange
	storeResponseFolder := "/tmp/responses"
	domain := "example.com"
	url := "http://example.com/path"
	expectedFileName := filepath.Join(storeResponseFolder, domain, getResponseHash(url)+".txt")

	// Act
	actualFileName := getResponseFileName(storeResponseFolder, domain, url)

	// Assert
	assert.Equal(t, expectedFileName, actualFileName, "File path should match the expected value")
}

// TestUpdateIndex_ValidInput_654 tests the `updateIndex` function with valid inputs.

func TestGetResponseHash_ValidURL_123(t *testing.T) {
	// Arrange
	url := "http://example.com"
	expectedHash := sha1.Sum([]byte(url))
	expectedHashString := hex.EncodeToString(expectedHash[:])

	// Act
	actualHash := getResponseHash(url)

	// Assert
	assert.Equal(t, expectedHashString, actualHash, "Hash should match the expected value")
}

// TestGetResponseHost_ValidURL_456 tests the `getResponseHost` function with a valid URL input.

func TestCreateHostDir_ValidInput_789(t *testing.T) {
	// Arrange
	storeResponseFolder := "/tmp/responses"
	domain := "example.com"

	// Act
	actualPath := createHostDir(storeResponseFolder, domain)

	// Assert
	expectedPath := filepath.Join(storeResponseFolder, domain)
	assert.Equal(t, expectedPath, actualPath, "Directory path should match the expected value")
	_, err := os.Stat(expectedPath)
	assert.NoError(t, err, "Directory should be created successfully")
}

// TestGetResponseFileName_ValidInput_321 tests the `getResponseFileName` function with valid inputs.

