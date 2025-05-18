package navigation

import (
	"testing"

	"net/http"
	"net/url"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHeaders_MarshalJSON_123 tests the MarshalJSON method of the Headers type.
func TestResponse_AbsoluteURL_456(t *testing.T) {
	// Arrange
	resp := &http.Response{
		Request: &http.Request{
			URL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/base",
			},
		},
	}
	response := Response{Resp: resp}

	// Test cases
	testCases := []struct {
		name     string
		path     string
		expected string
	}{
		{"Relative Path", "/test", "https://example.com/test"},
		{"Fragment Path", "#section", ""},
		{"Protocol-relative URL", "//other.com/path", "https://other.com/path"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := response.AbsoluteURL(tc.path)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestResponse_IsRedirect_789 tests the IsRedirect method of the Response struct.

func TestHeaders_MarshalJSON_123(t *testing.T) {
	// Arrange
	headers := Headers{
		"Content-Type":  "application/json",
		"Authorization": "Bearer token",
	}
	expectedJSON := `{"authorization":"Bearer token","content-type":"application/json"}`

	// Act
	marshaledJSON, err := headers.MarshalJSON()

	// Assert
	require.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(marshaledJSON))
}

// TestResponse_AbsoluteURL_456 tests the AbsoluteURL method of the Response struct.

func TestResponse_IsRedirect_789(t *testing.T) {
	// Test cases
	testCases := []struct {
		name       string
		statusCode int
		expected   bool
	}{
		{"Redirect Status Code 301", 301, true},
		{"Redirect Status Code 302", 302, true},
		{"Non-Redirect Status Code 200", 200, false},
		{"Non-Redirect Status Code 404", 404, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			response := Response{StatusCode: tc.statusCode}

			// Act
			result := response.IsRedirect()

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

