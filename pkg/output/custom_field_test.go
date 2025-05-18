package output

import (
	"regexp"
	"testing"

	"os"

	"path/filepath"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

// TestSetCompiledRegexp_AppendRegex_001 tests the behavior of `SetCompiledRegexp` method when appending a compiled regex.
func TestLoadCustomFields_ValidFile_005(t *testing.T) {
	// Arrange
	tempFile, err := os.CreateTemp("", "customfield*.yaml")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	yamlContent := `
     - name: valid_name
       type: regex
       part: response
       regex:
         - "test-regex"
     `
	_, err = tempFile.Write([]byte(yamlContent))
	require.NoError(t, err)

	// Act
	err = loadCustomFields(tempFile.Name(), "valid_name")

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, CustomFieldsMap, "valid_name")
}

// TestLoadCustomFields_InvalidRegex_007 tests the behavior of `loadCustomFields` function when the YAML file contains invalid regex patterns.

func TestParseCustomFieldName_ValidFile_003(t *testing.T) {
	// Arrange
	tempFile, err := os.CreateTemp("", "customfield*.yaml")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	yamlContent := `
     - name: valid_name
       type: regex
       part: response
       regex:
         - "test-regex"
     `
	_, err = tempFile.Write([]byte(yamlContent))
	require.NoError(t, err)

	// Act
	err = parseCustomFieldName(tempFile.Name())

	// Assert
	assert.NoError(t, err)
}

// TestParseCustomFieldName_InvalidName_004 tests the behavior of `parseCustomFieldName` function with invalid custom field names.

func TestInitCustomFieldConfigFile_Success_891(t *testing.T) {
	tempHomeDir := t.TempDir()

	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	require.NoError(t, os.Setenv("HOME", tempHomeDir))
	require.NoError(t, os.Setenv("USERPROFILE", tempHomeDir)) // For Windows compatibility
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("USERPROFILE", originalUserProfile)
	}()

	expectedConfigPath := filepath.Join(tempHomeDir, ".config", "katana", "field-config.yaml")

	path, err := initCustomFieldConfigFile()

	require.NoError(t, err)
	assert.Equal(t, expectedConfigPath, path)
	assert.FileExists(t, path)

	// Verify content is the default config
	content, readErr := os.ReadFile(path)
	require.NoError(t, readErr)
	var fields []CustomFieldConfig
	require.NoError(t, yaml.Unmarshal(content, &fields))
	assert.Equal(t, DefaultFieldConfigData, fields)
}

// TestInitCustomFieldConfigFile_FileExists_227 tests that if the field-config.yaml
// already exists, the function returns its path without modifying it.

func TestSetCompiledRegexp_AppendRegex_001(t *testing.T) {
	// Arrange
	regex := regexp.MustCompile(`test-regex`)
	config := &CustomFieldConfig{}

	// Act
	config.SetCompiledRegexp(regex)

	// Assert
	require.NotNil(t, config.CompileRegex)
	assert.Equal(t, 1, len(config.CompileRegex))
	assert.Equal(t, regex, config.CompileRegex[0])
}

// TestGetName_ReturnName_002 tests the behavior of `GetName` method to ensure it returns the correct name.

func TestParseCustomFieldName_InvalidName_004(t *testing.T) {
	// Arrange
	tempFile, err := os.CreateTemp("", "customfield*.yaml")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	yamlContent := `
     - name: invalid name
       type: regex
       part: response
       regex:
         - "test-regex"
     `
	_, err = tempFile.Write([]byte(yamlContent))
	require.NoError(t, err)

	// Act
	err = parseCustomFieldName(tempFile.Name())

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wrong custom field name")
}

// TestLoadCustomFields_ValidFile_005 tests the behavior of `loadCustomFields` function with a valid YAML file.

func TestParseCustomFieldName_OpenFileError_781(t *testing.T) {
	err := parseCustomFieldName("non_existent_file.yaml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not read field config")
	assert.Contains(t, err.Error(), "no such file or directory")
}

// TestLoadCustomFields_OpenFileError_826 tests that loadCustomFields
// returns an error when the specified file cannot be opened.

func TestInitCustomFieldConfigFile_FileExists_227(t *testing.T) {
	tempHomeDir := t.TempDir()

	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	require.NoError(t, os.Setenv("HOME", tempHomeDir))
	require.NoError(t, os.Setenv("USERPROFILE", tempHomeDir))
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("USERPROFILE", originalUserProfile)
	}()

	expectedConfigDir := filepath.Join(tempHomeDir, ".config", "katana")
	expectedConfigPath := filepath.Join(expectedConfigDir, "field-config.yaml")

	// Pre-create the directory and file
	require.NoError(t, os.MkdirAll(expectedConfigDir, 0755))
	file, err := os.Create(expectedConfigPath)
	require.NoError(t, err)
	_, err = file.WriteString("existing dummy content")
	require.NoError(t, err)
	require.NoError(t, file.Close())

	path, err := initCustomFieldConfigFile()

	require.NoError(t, err)
	assert.Equal(t, expectedConfigPath, path)

	// Verify content is unchanged
	content, readErr := os.ReadFile(path)
	require.NoError(t, readErr)
	assert.Equal(t, "existing dummy content", string(content))
}

// TestInitCustomFieldConfigFile_UserHomeDirError_538 attempts to test the scenario
// where os.UserHomeDir() returns an error. This is OS-dependent and might be skipped.

func TestInitCustomFieldConfigFile_MkdirError_991(t *testing.T) {
	tempHomeDir := t.TempDir()

	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	require.NoError(t, os.Setenv("HOME", tempHomeDir))
	require.NoError(t, os.Setenv("USERPROFILE", tempHomeDir))
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("USERPROFILE", originalUserProfile)
	}()

	// Create a file where ".config" directory is supposed to be, to cause MkdirAll to fail at ".config/katana"
	require.NoError(t, os.Mkdir(filepath.Join(tempHomeDir, ".config"), 0755)) // parent .config
	filePathAsDir := filepath.Join(tempHomeDir, ".config", "katana")          // This should be a dir
	file, err := os.Create(filePathAsDir)                                     // Create 'katana' as a file
	require.NoError(t, err)
	require.NoError(t, file.Close())

	path, err := initCustomFieldConfigFile()

	require.Error(t, err)
	// Error messages vary by OS, e.g., "not a directory" or "The system cannot find the path specified."
	// Check that the error message contains the problematic path segment.
	assert.Contains(t, err.Error(), filepath.Join(".config", "katana"))
	assert.Empty(t, path)
}

func TestGetName_ReturnName_002(t *testing.T) {
	// Arrange
	config := &CustomFieldConfig{Name: "test-name"}

	// Act
	name := config.GetName()

	// Assert
	assert.Equal(t, "test-name", name)
}

// TestParseCustomFieldName_ValidFile_003 tests the behavior of `parseCustomFieldName` function with a valid YAML file.

func TestLoadCustomFields_OpenFileError_826(t *testing.T) {
	CustomFieldsMap = make(map[string]CustomFieldConfig) // Reset
	defer func() { CustomFieldsMap = make(map[string]CustomFieldConfig) }()

	err := loadCustomFields("non_existent_file.yaml", "any_field")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not read field config")
	assert.Contains(t, err.Error(), "no such file or directory")
}

// TestInitCustomFieldConfigFile_Success_891 tests the successful creation of
// the field-config.yaml file when it does not already exist.
// It manipulates HOME/USERPROFILE environment variables to control the config path.

func TestLoadCustomFields_InvalidRegex_007(t *testing.T) {
	// Arrange
	tempFile, err := os.CreateTemp("", "customfield*.yaml")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	yamlContent := `
      - name: valid_name
        type: regex
        part: response
        regex:
          - "[invalid-regex("
      `
	_, err = tempFile.Write([]byte(yamlContent))
	require.NoError(t, err)

	// Act
	err = loadCustomFields(tempFile.Name(), "valid_name")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not parse regex in field config")
}

// TestParseCustomFieldName_OpenFileError_781 tests that parseCustomFieldName
// returns an error when the specified file cannot be opened.

func TestInitCustomFieldConfigFile_UserHomeDirError_538(t *testing.T) {
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	originalHomeDrive := os.Getenv("HOMEDRIVE")
	originalHomePath := os.Getenv("HOMEPATH")

	// Attempt to make os.UserHomeDir fail by unsetting typical env vars
	os.Unsetenv("HOME")
	os.Unsetenv("USERPROFILE")
	os.Unsetenv("HOMEDRIVE")
	os.Unsetenv("HOMEPATH")

	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("USERPROFILE", originalUserProfile)
		os.Setenv("HOMEDRIVE", originalHomeDrive)
		os.Setenv("HOMEPATH", originalHomePath)
	}()

	// Check if os.UserHomeDir actually fails under these conditions
	_, homeErr := os.UserHomeDir()
	if homeErr == nil {
		t.Skip("Skipping UserHomeDir error test: os.UserHomeDir did not fail with unset HOME/USERPROFILE like env vars. This path may be untestable in the current environment without source modification.")
	}

	path, err := initCustomFieldConfigFile()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not get home directory")
	assert.Empty(t, path)
}

// TestInitCustomFieldConfigFile_MkdirError_991 tests the scenario where
// creating the .config/katana directory fails.

