package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAddSchemeIfNotExists_WithScheme_004 tests the behavior of addSchemeIfNotExists when the input URL already has a scheme.
func TestAddSchemeIfNotExists_WithoutScheme_005(t *testing.T) {
	inputURL := "example.com"
	result := addSchemeIfNotExists(inputURL)
	assert.Equal(t, "https://example.com", result)
}

// TestAddSchemeIfNotExists_InvalidURL_006 tests the behavior of addSchemeIfNotExists when the input URL is invalid.

func TestAddSchemeIfNotExists_WithScheme_004(t *testing.T) {
	inputURL := "https://example.com"
	result := addSchemeIfNotExists(inputURL)
	assert.Equal(t, inputURL, result)
}

// TestAddSchemeIfNotExists_WithoutScheme_005 tests the behavior of addSchemeIfNotExists when the input URL does not have a scheme.

func TestAddSchemeIfNotExists_InvalidURL_006(t *testing.T) {
	inputURL := "://invalid-url"
	result := addSchemeIfNotExists(inputURL)
	assert.Equal(t, inputURL, result)
}

// TestAddSchemeIfNotExists_HTTPPort80_707 tests addSchemeIfNotExists for an input URL
// with port 80, expecting it to be prefixed with http://.

func TestAddSchemeIfNotExists_HTTPPort80_707(t *testing.T) {
	inputURL := "example.com:80"
	expectedURL := "http://example.com:80"
	result := addSchemeIfNotExists(inputURL)
	assert.Equal(t, expectedURL, result)
}

