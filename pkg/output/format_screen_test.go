package output

import (
	"testing"

	"github.com/logrusorgru/aurora"
	"github.com/projectdiscovery/katana/pkg/navigation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFormatScreen_FieldsEmptyVerboseFalse_003 tests the behavior of formatScreen when `fields` is empty and `verbose` is false.
func TestFormatScreen_NoFields_VerboseTrue_AllParts_555(t *testing.T) {
	// Arrange
	au := aurora.NewAurora(false)
	writer := &StandardWriter{
		fields:  "",
		verbose: true,
		aurora:  au,
	}
	result := &Result{
		Request: &navigation.Request{
			Tag:    "test-tag",
			Method: "POST",
			URL:    "http://example.com/submit",
			Body:   "payload",
		},
	}
	expectedOutput := "[test-tag] [POST] http://example.com/submit [payload]"

	// Act
	output, err := writer.formatScreen(result)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, string(output))
}

// TestFormatScreen_WithUnknownField_VerboseTrue_777 tests formatting when an unknown field is specified and verbose is true.
// This ensures that if formatField returns no results, the output is empty.

func TestFormatScreen_FieldsEmptyVerboseFalse_003(t *testing.T) {
	// Arrange
	writer := &StandardWriter{
		fields:  "",
		verbose: false,
	}

	result := &Result{
		Request: &navigation.Request{
			URL: "http://example.com",
		},
	}

	// Act
	output, err := writer.formatScreen(result)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "http://example.com", string(output))
}

// TestFormatScreen_NoFields_VerboseTrue_AllParts_555 tests formatting when no fields are specified, verbose is true, and all optional request parts are present.

func TestFormatScreen_SingleField_URL_VerboseTrue_501(t *testing.T) {
	// Arrange
	au := aurora.NewAurora(false)
	writer := &StandardWriter{
		fields:  "url",
		verbose: true,
		aurora:  au,
	}
	testURL := "http://example.com/specific-url"
	result := &Result{
		Request: &navigation.Request{
			URL:    testURL,
			Tag:    "tag1", // Other data to ensure they are ignored
			Method: "POST",
		},
	}
	// Assuming formatField for "url" returns [{field: "url", value: testURL}]
	expectedOutput := "[url] " + testURL + "\n"

	// Act
	output, err := writer.formatScreen(result)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, string(output))
}

func TestFormatScreen_WithUnknownField_VerboseTrue_777(t *testing.T) {
	// Arrange
	au := aurora.NewAurora(false)
	writer := &StandardWriter{
		fields:  "nonexistentfield", // A field that formatField won't find
		verbose: true,
		aurora:  au,
	}
	result := &Result{
		Request: &navigation.Request{
			URL: "http://example.com", // Provide a base request
		},
	}
	expectedOutput := "" // Expect empty output as the field loop won't run

	// Act
	output, err := writer.formatScreen(result)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, string(output))
}

// TestFormatScreen_SingleField_URL_VerboseTrue_501 verifies output for a single specified field ("url") with verbose mode.

