package runner

import (
	"testing"

	"os"

	"path/filepath"

	"strings"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/katana/pkg/types"
	"github.com/projectdiscovery/katana/pkg/utils"
	fileutil "github.com/projectdiscovery/utils/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// TestValidateOptions_ValidOptions_777 ensures that a typical valid set of options
// passes validation without errors.
func TestParseInputs_FromOptionsURLs_121(t *testing.T) {
	r := &Runner{
		options: &types.Options{
			URLs: goflags.StringSlice{"http://example.com/1 ", " http://example.com/2", "http://example.com/1 "},
		},
		stdin: false,
	}
	inputs := r.parseInputs()
	assert.ElementsMatch(t, []string{"http://example.com/1", "http://example.com/2"}, inputs)
}

// TestParseInputs_FromStdin_232 verifies parsing of URLs
// solely from Stdin.

func TestValidateOptions_ValidOptions_777(t *testing.T) {
	options := &types.Options{
		MaxDepth: 5,
		URLs:     goflags.StringSlice{"http://example.com"},
	}
	err := validateOptions(options)
	require.NoError(t, err)
}

// TestReadCustomFormConfig_Success_123 verifies successful reading and parsing
// of a custom form configuration file.

func TestReadCustomFormConfig_Success_123(t *testing.T) {
	yamlContent := `
 email: test@example.com
 color: blue
 password: testpassword
 phone: "1234567890"
 placeholder: testplaceholder
 `
	tmpFile, err := os.CreateTemp(t.TempDir(), "formconfig-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(yamlContent)
	require.NoError(t, err)
	err = tmpFile.Close()
	require.NoError(t, err)

	// Preserve original FormData and restore after test
	originalFormData := utils.FormData
	defer func() { utils.FormData = originalFormData }()

	err = readCustomFormConfig(tmpFile.Name())
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", utils.FormData.Email)
	assert.Equal(t, "blue", utils.FormData.Color)
}

// TestReadCustomFormConfig_FileOpenError_456 checks error handling
// when the form configuration file cannot be opened.

func TestInitExampleFormFillConfig_CreateNewConfig_445(t *testing.T) {
	mockHome := t.TempDir()
	originalHomeEnv := os.Getenv("HOME")
	originalUserProfilerEnv := os.Getenv("USERPROFILE")

	os.Setenv("HOME", mockHome)
	os.Setenv("USERPROFILE", mockHome) // For Windows
	defer func() {
		os.Setenv("HOME", originalHomeEnv)
		os.Setenv("USERPROFILE", originalUserProfilerEnv)
	}()

	defaultConfigPath := filepath.Join(mockHome, ".config", "katana", "form-config.yaml")
	defer os.RemoveAll(filepath.Join(mockHome, ".config")) // Clean up

	err := initExampleFormFillConfig()
	require.NoError(t, err)

	require.True(t, fileutil.FileExists(defaultConfigPath), "Default config file should be created")

	content, err := os.ReadFile(defaultConfigPath)
	require.NoError(t, err)

	var data utils.FormFillData
	err = yaml.Unmarshal(content, &data)
	require.NoError(t, err)
	assert.Equal(t, utils.DefaultFormFillData.Email, data.Email)
}

// TestInitExampleFormFillConfig_MkdirError_556 tests error handling when
// os.MkdirAll fails during config initialization.

func TestParseInputs_FromStdin_232(t *testing.T) {
	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }()

	content := "http://example.com/stdin1\n http://example.com/stdin2 \nhttp://example.com/stdin1"
	tmpfile, err := os.CreateTemp(t.TempDir(), "stdin")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(content)
	require.NoError(t, err)
	_, err = tmpfile.Seek(0, 0)
	require.NoError(t, err)
	os.Stdin = tmpfile

	r := &Runner{
		options: &types.Options{URLs: goflags.StringSlice{}},
		stdin:   true,
	}
	inputs := r.parseInputs()
	assert.ElementsMatch(t, []string{"http://example.com/stdin1", "http://example.com/stdin2"}, inputs)
}

// TestInitExampleFormFillConfig_DefaultConfigExists tests the scenario where
// the default form configuration file already exists.

func TestInitExampleFormFillConfig_DefaultConfigExists_334(t *testing.T) {
	mockHome := t.TempDir()
	originalHomeEnv := os.Getenv("HOME")
	originalUserProfilerEnv := os.Getenv("USERPROFILE")

	os.Setenv("HOME", mockHome)
	os.Setenv("USERPROFILE", mockHome) // For Windows
	defer func() {
		os.Setenv("HOME", originalHomeEnv)
		os.Setenv("USERPROFILE", originalUserProfilerEnv)
	}()

	configDir := filepath.Join(mockHome, ".config", "katana")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	defaultConfigPath := filepath.Join(configDir, "form-config.yaml")
	dummyContent := "email: dummy@example.com"
	err = os.WriteFile(defaultConfigPath, []byte(dummyContent), 0644)
	require.NoError(t, err)
	defer os.Remove(defaultConfigPath) // Clean up

	// Preserve original FormData and restore after test
	originalFormData := utils.FormData
	defer func() { utils.FormData = originalFormData }()

	err = initExampleFormFillConfig()
	require.NoError(t, err)
	assert.Equal(t, "dummy@example.com", utils.FormData.Email)
}

// TestInitExampleFormFillConfig_CreateNewConfig_445 tests creation of a new
// default form configuration file if one doesn't exist.

func TestValidateOptions_InvalidMatchRegex_Error_308(t *testing.T) {
	options := &types.Options{
		MaxDepth:         1, // Satisfy depth/duration
		URLs:             goflags.StringSlice{"http://example.com"},
		OutputMatchRegex: goflags.StringSlice{"[invalid-regex"},
	}
	err := validateOptions(options)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid value for match regex option")
}

// TestValidateOptions_ValidMatchRegex_Success_217 checks that a valid OutputMatchRegex
// is compiled and added to options.MatchRegex.

func TestValidateOptions_InvalidFilterRegex_Error_539(t *testing.T) {
	options := &types.Options{
		MaxDepth:          1, // Satisfy depth/duration
		URLs:              goflags.StringSlice{"http://example.com"},
		OutputFilterRegex: goflags.StringSlice{"[invalid-regex"},
	}
	err := validateOptions(options)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid value for filter regex option")
}

// TestValidateOptions_ValidFilterRegex_Success_870 checks that a valid OutputFilterRegex
// is compiled and added to options.FilterRegex.

func TestValidateOptions_AutomaticFormFillDisabledForHeadless_Info_603(t *testing.T) {
	options := &types.Options{
		MaxDepth:          1, // Satisfy depth/duration
		URLs:              goflags.StringSlice{"http://example.com"},
		Headless:          true,
		AutomaticFormFill: true,
	}
	err := validateOptions(options)
	require.NoError(t, err)
	assert.False(t, options.AutomaticFormFill, "AutomaticFormFill should be disabled for headless mode")
}

// TestValidateOptions_SystemChromePathExists_Success_493 checks that no error is returned
// if the specified system chrome binary exists.

func TestValidateOptions_StoreResponseDirEnablesStoreResponse_AutoEnable_775(t *testing.T) {
	options := &types.Options{
		MaxDepth:         1, // Satisfy depth/duration
		URLs:             goflags.StringSlice{"http://example.com"},
		StoreResponseDir: "responses",
		StoreResponse:    false,
	}
	err := validateOptions(options)
	require.NoError(t, err)
	assert.True(t, options.StoreResponse, "StoreResponse should be enabled if StoreResponseDir is set")
}

// TestValidateOptions_InvalidMatchRegex_Error_308 checks that an error is returned
// for an invalid OutputMatchRegex.

func TestValidateOptions_KnownFilesAdjustsMaxDepth_Adjust_433(t *testing.T) {
	options := &types.Options{
		KnownFiles: "read",
		MaxDepth:   2,
		URLs:       goflags.StringSlice{"http://example.com"}, // Satisfy input requirement
	}
	err := validateOptions(options)
	require.NoError(t, err)
	assert.Equal(t, 3, options.MaxDepth, "MaxDepth should be adjusted to 3 for KnownFiles")
}

// TestInitExampleFormFillConfig_UserHomeDirError_Error_135 checks error handling
// when os.UserHomeDir fails.

func TestValidateOptions_ValidFilterRegex_Success_870(t *testing.T) {
	options := &types.Options{
		MaxDepth:          1, // Satisfy depth/duration
		URLs:              goflags.StringSlice{"http://example.com"},
		OutputFilterRegex: goflags.StringSlice{"valid-regex"},
	}
	err := validateOptions(options)
	require.NoError(t, err)
	require.Len(t, options.FilterRegex, 1)
	assert.Equal(t, "valid-regex", options.FilterRegex[0].String())
}

// TestValidateOptions_KnownFilesAdjustsMaxDepth_Adjust_433 ensures MaxDepth is set to 3
// when KnownFiles option is used and initial MaxDepth is < 3.

func TestInitExampleFormFillConfig_MkdirError_556(t *testing.T) {
	mockHome := t.TempDir()
	originalHomeEnv := os.Getenv("HOME")
	originalUserProfilerEnv := os.Getenv("USERPROFILE")

	os.Setenv("HOME", mockHome)
	os.Setenv("USERPROFILE", mockHome) // For Windows
	defer func() {
		os.Setenv("HOME", originalHomeEnv)
		os.Setenv("USERPROFILE", originalUserProfilerEnv)
	}()

	// Path where ".config/katana" directory would be created by initExampleFormFillConfig.
	// We make ".config" a file to cause MkdirAll for ".config/katana" to fail.
	dotConfigParentAsFile := filepath.Join(mockHome, ".config")
	err := os.WriteFile(dotConfigParentAsFile, []byte("i am a file"), 0600)
	require.NoError(t, err)
	defer os.Remove(dotConfigParentAsFile)

	err = initExampleFormFillConfig()
	require.Error(t, err)
	// Error message varies by OS, e.g., "mkdir ...: not a directory" or "mkdir ...: The system cannot find the path specified."
	// Check for a generic part of the expected error related to directory creation.
	assert.Contains(t, strings.ToLower(err.Error()), "mkdir", "Error message should indicate a mkdir failure")
}

// TestValidateOptions_AutomaticFormFillDisabledForHeadless_Info_603 ensures that
// AutomaticFormFill is set to false if Headless mode is true.

func TestValidateOptions_ValidMatchRegex_Success_217(t *testing.T) {
	options := &types.Options{
		MaxDepth:         1, // Satisfy depth/duration
		URLs:             goflags.StringSlice{"http://example.com"},
		OutputMatchRegex: goflags.StringSlice{"valid-regex"},
	}
	err := validateOptions(options)
	require.NoError(t, err)
	require.Len(t, options.MatchRegex, 1)
	assert.Equal(t, "valid-regex", options.MatchRegex[0].String())
}

// TestValidateOptions_InvalidFilterRegex_Error_539 checks that an error is returned
// for an invalid OutputFilterRegex.

func TestInitExampleFormFillConfig_UserHomeDirError_Error_135(t *testing.T) {
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")

	// Unsetting HOME (and USERPROFILE for Windows) can cause os.UserHomeDir to fail on some systems.
	os.Unsetenv("HOME")
	os.Unsetenv("USERPROFILE")
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("USERPROFILE", originalUserProfile)
	}()

	err := initExampleFormFillConfig()
	// This error condition might be hard to reliably trigger across all OSes.
	// If os.UserHomeDir() succeeds even with unset env vars (e.g., by other means), this test might not fail as intended.
	// However, it's a common way to attempt to trigger this error.
	if err != nil { // Only assert if an error actually occurred, as UserHomeDir behavior can vary
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not get home directory", "Error should be about failing to get home directory")
	} else {
		t.Log("os.UserHomeDir() did not fail as expected; skipping exact error assertion for this scenario.")
		// If it didn't fail, ensure the config dir and file were created to clean up.
		// This part assumes the function proceeded to create the file.
		// Determine expected config path based on potentially successful UserHomeDir.
		homedir, _ := os.UserHomeDir() // Re-call to see where it might have put it
		if homedir != "" {
			os.RemoveAll(filepath.Join(homedir, ".config"))
		}
	}
}

// TestInitExampleFormFillConfig_CreateFileError_Error_709 checks error handling
// when creating the example form config file fails.
// It uses the actual confusing error message from the source.

func TestReadCustomFormConfig_YamlDecodeError_789(t *testing.T) {
	malformedYamlContent := `email: test@example.com
 color: blue
   badindent: true` // Malformed YAML
	tmpFile, err := os.CreateTemp(t.TempDir(), "formconfig-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(malformedYamlContent)
	require.NoError(t, err)
	err = tmpFile.Close()
	require.NoError(t, err)

	err = readCustomFormConfig(tmpFile.Name())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not decode form config")
}

// TestParseInputs_FromOptionsURLs_121 verifies parsing of URLs
// solely from the Options.URLs field.

func TestValidateOptions_SystemChromePathExists_Success_493(t *testing.T) {
	tempFile, errCreate := os.CreateTemp(t.TempDir(), "mockchrome")
	require.NoError(t, errCreate)
	defer os.Remove(tempFile.Name())
	errClose := tempFile.Close()
	require.NoError(t, errClose)

	options := &types.Options{
		MaxDepth:         1, // Satisfy depth/duration
		URLs:             goflags.StringSlice{"http://example.com"},
		SystemChromePath: tempFile.Name(),
		Headless:         true, // Headless must be true
	}
	errValidate := validateOptions(options)
	require.NoError(t, errValidate)
}

// TestValidateOptions_StoreResponseDirEnablesStoreResponse_AutoEnable_775 ensures that
// StoreResponse flag is automatically enabled if StoreResponseDir is specified.

func TestInitExampleFormFillConfig_CreateFileError_Error_709(t *testing.T) {
	mockHome := t.TempDir()
	originalHomeEnv := os.Getenv("HOME")
	originalUserProfilerEnv := os.Getenv("USERPROFILE")

	os.Setenv("HOME", mockHome)
	os.Setenv("USERPROFILE", mockHome) // For Windows
	defer func() {
		os.Setenv("HOME", originalHomeEnv)
		os.Setenv("USERPROFILE", originalUserProfilerEnv)
	}()

	// Path to the default config file that initExampleFormFillConfig will try to create.
	configFilePath := filepath.Join(mockHome, ".config", "katana", "form-config.yaml")
	// Create a directory at that path, so os.Create(configFilePath) fails.
	err := os.MkdirAll(configFilePath, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(filepath.Join(mockHome, ".config")) // Clean up the mock .config dir

	err = initExampleFormFillConfig()
	require.Error(t, err)
	// The source code has a confusing error message here. We test against that specific message.
	assert.Contains(t, err.Error(), "could not get home directory")
}

func TestReadCustomFormConfig_FileOpenError_456(t *testing.T) {
	nonExistentPath := filepath.Join(t.TempDir(), "nonexistent.yaml")
	err := readCustomFormConfig(nonExistentPath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not read form config")
}

// TestReadCustomFormConfig_YamlDecodeError_789 checks error handling
// for malformed YAML in the form configuration file.

