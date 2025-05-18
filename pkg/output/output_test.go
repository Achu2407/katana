package output

import (
	"testing"

	"os"
	"path/filepath"
	"regexp"

	"strings"
	"sync"
	"time"

	"errors"

	jsoniter "github.com/json-iterator/go"
	"github.com/logrusorgru/aurora"
	"github.com/projectdiscovery/katana/pkg/navigation"
	"github.com/projectdiscovery/katana/pkg/utils/extensions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWrite_NilResult_002 tests the Write method when the result is nil.
func TestNew_ValidOptions_Fixed_781(t *testing.T) {
	originalStoreFieldDir := storeFieldDir
	defer func() { storeFieldDir = originalStoreFieldDir }()

	tempDir, err := os.MkdirTemp("", "katana-test-new")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	fieldConfigContent := `
 - name: field1
   type: regex
   part: body
   regex: [".*"]
 - name: field2
   type: regex
   part: body
   regex: [".*"]
 - name: storeField1
   type: regex
   part: body
   regex: [".*"]
 - name: storeField2
   type: regex
   part: body
   regex: [".*"]
 `
	tempFieldConfigFile := filepath.Join(tempDir, "config.yaml")
	err = os.WriteFile(tempFieldConfigFile, []byte(fieldConfigContent), 0644)
	require.NoError(t, err)

	options := Options{
		Colors:                true,
		JSON:                  true,
		Verbose:               true,
		StoreResponse:         true,
		NoClobber:             false,
		OmitRaw:               false,
		OmitBody:              false,
		OutputFile:            filepath.Join(tempDir, "output.txt"),
		Fields:                "field1,field2",
		StoreFields:           "storeField1,storeField2",
		StoreResponseDir:      filepath.Join(tempDir, "responseDir"),
		StoreFieldDir:         filepath.Join(tempDir, "fieldDir"),
		FieldConfig:           tempFieldConfigFile,
		ErrorLogFile:          filepath.Join(tempDir, "error.log"),
		MatchRegex:            []*regexp.Regexp{regexp.MustCompile(".*")},
		FilterRegex:           []*regexp.Regexp{regexp.MustCompile("nonexistent")},
		ExtensionValidator:    extensions.NewValidator(nil, nil),
		OutputMatchCondition:  "",
		OutputFilterCondition: "",
	}

	writer, err := New(options)
	require.NoError(t, err)
	require.NotNil(t, writer)

	standardWriter, ok := writer.(*StandardWriter)
	require.True(t, ok)

	assert.Equal(t, options.Fields, standardWriter.fields)
	assert.Equal(t, options.JSON, standardWriter.json)
	assert.Equal(t, options.Verbose, standardWriter.verbose)
	assert.Equal(t, options.StoreResponse, standardWriter.storeResponse)
	assert.Equal(t, options.StoreResponseDir, standardWriter.storeResponseDir)
	assert.Equal(t, options.NoClobber, standardWriter.noClobber)
	assert.Equal(t, options.OmitRaw, standardWriter.omitRaw)
	assert.Equal(t, options.OmitBody, standardWriter.omitBody)
	assert.Equal(t, options.MatchRegex, standardWriter.matchRegex)
	assert.Equal(t, options.FilterRegex, standardWriter.filterRegex)
	assert.Equal(t, options.ExtensionValidator, standardWriter.extensionValidator)
	assert.Equal(t, options.OutputMatchCondition, standardWriter.outputMatchCondition)
	assert.Equal(t, options.OutputFilterCondition, standardWriter.outputFilterCondition)

	assert.NotNil(t, standardWriter.outputFile)
	assert.NotNil(t, standardWriter.errorFile)
	assert.DirExists(t, standardWriter.storeResponseDir)
	assert.FileExists(t, filepath.Join(standardWriter.storeResponseDir, indexFile))

	if options.StoreFields != "" {
		expectedStoreFieldDir := options.StoreFieldDir
		if options.StoreFieldDir == "" { // Should not happen with this test's options
			expectedStoreFieldDir = originalStoreFieldDir
		}
		assert.DirExists(t, expectedStoreFieldDir)
		assert.Equal(t, expectedStoreFieldDir, storeFieldDir) // global var check
	}
}

// TestWriteErr_ValidErrorMessage_Fixed_112 tests that WriteErr correctly writes
// an error message to the error log file. It uses a temporary file to capture the output.

func TestEvalDslExpr_TrueCondition_703(t *testing.T) {
	result := &Result{
		Response: &navigation.Response{StatusCode: 200},
	}
	dslExpr := "status_code == 200" // 'status_code' comes from flattening Response.StatusCode
	assert.True(t, evalDslExpr(result, dslExpr))
}

// TestIgnoreErr_ShowDSLErrFalse_101 tests ignoreErr with showDSLErr = false

func TestNew_StoreResponse_NoClobber_CreateNewDir_303(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "katana-test-noclobber")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a dummy field config
	fieldConfigContent := `[]`
	tempFieldConfigFile := filepath.Join(tempDir, "config.yaml")
	err = os.WriteFile(tempFieldConfigFile, []byte(fieldConfigContent), 0644)
	require.NoError(t, err)

	baseResponseDirName := "responses"
	responseDirParent := tempDir
	responseDir := filepath.Join(responseDirParent, baseResponseDirName)

	// Create the initial directory to trigger NoClobber logic
	err = os.MkdirAll(responseDir, 0755)
	require.NoError(t, err)

	options := Options{
		StoreResponse:    true,
		NoClobber:        true,
		StoreResponseDir: responseDir,         // Target the existing directory
		FieldConfig:      tempFieldConfigFile, // Provide a minimal valid config
	}

	writer, err := New(options)
	require.NoError(t, err)
	require.NotNil(t, writer)

	standardWriter, ok := writer.(*StandardWriter)
	require.True(t, ok)

	// Expect a new directory like "responses1"
	expectedNewResponseDir := filepath.Join(responseDirParent, baseResponseDirName+"1")
	assert.Equal(t, expectedNewResponseDir, standardWriter.storeResponseDir)
	assert.DirExists(t, standardWriter.storeResponseDir)
	assert.FileExists(t, filepath.Join(standardWriter.storeResponseDir, indexFile))
}

// TestWrite_ExtensionFilteredOut_401 verifies that Write returns an error
// if the result URL is filtered out by the extension validator.

func TestWriteErr_ValidErrorMessage_Fixed_112(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "katana-test-writeerr")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tempErrorLogFile := filepath.Join(tempDir, "error.log")

	errorFw, err := newFileOutputWriter(tempErrorLogFile)
	require.NoError(t, err)
	// No need to defer errorFw.Close() here as StandardWriter.Close() would handle it,
	// but this test doesn't call StandardWriter.Close(). We'll close it manually if needed for read.

	writer := &StandardWriter{
		errorFile:   errorFw,
		outputMutex: &sync.Mutex{},
	}

	errMsg := &Error{
		Timestamp: time.Now().UTC().Truncate(time.Second), // Truncate for consistent comparison
		Endpoint:  "http://example.com/error",
		Source:    "test_source",
		Error:     "a test error occurred",
	}

	err = writer.WriteErr(errMsg)
	require.NoError(t, err)

	// Close the file writer to ensure data is flushed before reading
	err = errorFw.Close()
	require.NoError(t, err)

	content, err := os.ReadFile(tempErrorLogFile)
	require.NoError(t, err)

	expectedJSON, err := jsoniter.Marshal(errMsg)
	require.NoError(t, err)

	// newFileOutputWriter might add a newline, so we compare trimmed JSON strings
	assert.JSONEq(t, string(expectedJSON), strings.TrimSpace(string(content)))
}

// TestNew_StoreResponse_NoClobber_CreateNewDir_303 ensures that if StoreResponse and NoClobber are true,
// and the response directory already exists, a new directory with an incremented suffix is created.

func TestFilterOutput_RegexMatch_601(t *testing.T) {
	filterRegex := regexp.MustCompile(`.*\.css`)
	writer := &StandardWriter{
		filterRegex: []*regexp.Regexp{filterRegex},
	}
	event := &Result{
		Request: &navigation.Request{URL: "http://example.com/style.css"},
	}
	assert.True(t, writer.filterOutput(event))

	eventNoMatch := &Result{
		Request: &navigation.Request{URL: "http://example.com/index.html"},
	}
	assert.False(t, writer.filterOutput(eventNoMatch))
}

// TestEvalDslExpr_TrueCondition_703 checks if evalDslExpr correctly evaluates
// a DSL expression to true based on the provided result.

func TestStandardWriter_Close_FileErrors_405(t *testing.T) {
	tempDir := t.TempDir()

	// Helper to create a fileWriter whose underlying file is already closed,
	// so its Close() method (which calls Sync then Close on os.File) might error.
	createFailingFileWriter := func(path string) *fileWriter {
		fw, err := newFileOutputWriter(path)
		require.NoError(t, err)
		// Close the underlying os.File. Subsequent Sync or Close in fw.Close() should fail.
		err = fw.file.Close()   // Close the actual os.File
		require.NoError(t, err) // First close should be fine
		return fw
	}

	// Scenario 1: outputFile.Close() returns an error
	outputFilePath := filepath.Join(tempDir, "out.txt")
	failingOutputFw := createFailingFileWriter(outputFilePath)
	writerWithOutputError := &StandardWriter{outputFile: failingOutputFw}
	err := writerWithOutputError.Close()
	require.Error(t, err, "Expected error from outputFile.Close()")
	// The specific error can be OS-dependent (e.g., "invalid argument" or "bad file descriptor" from sync/close on closed fd)
	// contains check is safer.
	// Error often contains "bad file descriptor" or "invalid argument" when operating on closed file.
	// Let's check for "file" as a generic term.
	assert.Contains(t, strings.ToLower(err.Error()), "file", "Error message should indicate a file operation issue")

	// Scenario 2: errorFile.Close() returns an error
	errorFilePath := filepath.Join(tempDir, "err.log")
	failingErrorFw := createFailingFileWriter(errorFilePath)
	writerWithErrorFileError := &StandardWriter{errorFile: failingErrorFw}
	err = writerWithErrorFileError.Close()
	require.Error(t, err, "Expected error from errorFile.Close()")
	assert.Contains(t, strings.ToLower(err.Error()), "file", "Error message should indicate a file operation issue")

	// Scenario 3: outputFile.Close() succeeds, errorFile.Close() fails
	stableOutputFilePath := filepath.Join(tempDir, "stable_out.txt")
	stableOutputFw, _ := newFileOutputWriter(stableOutputFilePath) // Normal file writer

	failingErrorFw2 := createFailingFileWriter(filepath.Join(tempDir, "err2.log"))
	writerWithErrorFileErrorOnly := &StandardWriter{
		outputFile: stableOutputFw,
		errorFile:  failingErrorFw2,
	}
	err = writerWithErrorFileErrorOnly.Close()
	require.Error(t, err, "Expected error from errorFile.Close() when outputFile is fine")
	assert.Contains(t, strings.ToLower(err.Error()), "file", "Error message should indicate a file operation issue from errorFile")
}

func TestStandardWriter_MatchAndFilterOutput_WithDSL_280(t *testing.T) {
	result200 := &Result{
		Request:  &navigation.Request{URL: "http://example.com/page"},
		Response: &navigation.Response{StatusCode: 200},
	}
	result404 := &Result{
		Request:  &navigation.Request{URL: "http://example.com/notfound"},
		Response: &navigation.Response{StatusCode: 404},
	}

	// Match DSL: status_code == 200
	writerMatchDSL := &StandardWriter{outputMatchCondition: "status_code == 200"}
	assert.True(t, writerMatchDSL.matchOutput(result200), "Should match result with status 200")
	assert.False(t, writerMatchDSL.matchOutput(result404), "Should not match result with status 404")

	// Filter DSL: status_code == 404
	writerFilterDSL := &StandardWriter{outputFilterCondition: "status_code == 404"}
	assert.True(t, writerFilterDSL.filterOutput(result404), "Should filter result with status 404")
	assert.False(t, writerFilterDSL.filterOutput(result200), "Should not filter result with status 200")

	// Both regex and DSL (matchOutput: regex OR dsl)
	writerMatchRegexAndDSL := &StandardWriter{
		matchRegex:           []*regexp.Regexp{regexp.MustCompile("otherdomain.com")}, // Won't match
		outputMatchCondition: "status_code == 200",                                    // Will match for result200
	}
	assert.True(t, writerMatchRegexAndDSL.matchOutput(result200), "Should match via DSL if regex fails")
	assert.False(t, writerMatchRegexAndDSL.matchOutput(result404), "Should not match if both regex and DSL fail")

	// Both regex and DSL (filterOutput: regex OR dsl)
	writerFilterRegexAndDSL := &StandardWriter{
		filterRegex:           []*regexp.Regexp{regexp.MustCompile("otherdomain.com")}, // Won't match
		outputFilterCondition: "status_code == 404",                                    // Will match for result404
	}
	assert.True(t, writerFilterRegexAndDSL.filterOutput(result404), "Should filter via DSL if regex fails to filter")
	assert.False(t, writerFilterRegexAndDSL.filterOutput(result200), "Should not filter if both regex and DSL don't filter")
}

// TestEvalDslExpr_ErrorScenarios_290 tests DSL evaluation error paths.

func TestIgnoreErr_ShowDSLErrFalse_101(t *testing.T) {
	originalShowDSLErr := showDSLErr
	showDSLErr = false // Explicitly set for this test
	defer func() { showDSLErr = originalShowDSLErr }()

	noParamErr := errors.New("something No parameter something")
	assert.True(t, ignoreErr(noParamErr), "Error containing 'No parameter' should be ignored when showDSLErr is false")

	// Simulate dsl.ErrParsingArg using errors.New for testing the errors.Is branch (if dsl.ErrParsingArg were non-nil and matched)
	// For this test, we rely on the string check mostly.
	// If dsl.ErrParsingArg was `errors.New("specific parsing error")`, then:
	// dslParsingArgError := dsl.ErrParsingArg // if accessible, otherwise simulate
	// assert.True(t, ignoreErr(dslParsingArgError), "dsl.ErrParsingArg should be ignored when showDSLErr is false")

	standardErr := errors.New("a standard error")
	assert.False(t, ignoreErr(standardErr), "Standard error should not be ignored when showDSLErr is false")
}

// TestIgnoreErr_ShowDSLErrTrue_102 tests ignoreErr with showDSLErr = true

func TestClose_NilFiles_004(t *testing.T) {
	writer := &StandardWriter{}
	err := writer.Close()
	require.NoError(t, err)
}

// TestMatchOutput_NoConditions_005 tests the matchOutput method when no matchRegex or outputMatchCondition is provided.

func TestWrite_ExtensionFilteredOut_401(t *testing.T) {
	validator := extensions.NewValidator(nil, []string{".js"}) // Filter out .js files
	writer := &StandardWriter{
		extensionValidator: validator,
		aurora:             aurora.NewAurora(false), // For NoColor
		outputMutex:        &sync.Mutex{},
	}

	result := &Result{
		Request: &navigation.Request{
			URL: "http://example.com/script.js",
		},
	}

	err := writer.Write(result)
	require.Error(t, err)
	assert.Equal(t, "result does not match extension filter", err.Error())
}

// TestFilterOutput_RegexMatch_601 verifies that filterOutput returns true
// if the event URL matches one of the filter regexes.

func TestWrite_NilResult_002(t *testing.T) {
	writer := &StandardWriter{}
	err := writer.Write(nil)
	require.Error(t, err)
	assert.Equal(t, "result is nil", err.Error())
}

// TestClose_NilFiles_004 tests the Close method when both outputFile and errorFile are nil.

func TestMatchOutput_NoConditions_005(t *testing.T) {
	writer := &StandardWriter{}
	result := &Result{
		Request: &navigation.Request{
			URL: "http://example.com",
		},
	}

	match := writer.matchOutput(result)
	assert.True(t, match)
}

// TestNew_ValidOptions_Fixed_781 tests the New function with a comprehensive set of valid options.
// It ensures that the StandardWriter is initialized correctly and that file system operations
// (like creating output files and directories) are performed as expected.
// Temporary directories and files are used to isolate the test and prevent side effects.

func TestRemoveDirsWithSuffix_270(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "katana-test-remove-suffix")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	baseName := "data"
	dirToRemove := filepath.Join(tempDir, baseName)        // data
	dirToRemove1 := filepath.Join(tempDir, baseName+"1")   // data1
	dirToRemove10 := filepath.Join(tempDir, baseName+"10") // data10
	dirToKeep := filepath.Join(tempDir, baseName+"abc")    // dataabc
	fileToKeep := filepath.Join(tempDir, baseName+"_file.txt")

	require.NoError(t, os.Mkdir(dirToRemove, 0755))
	require.NoError(t, os.Mkdir(dirToRemove1, 0755))
	require.NoError(t, os.Mkdir(dirToRemove10, 0755))
	require.NoError(t, os.Mkdir(dirToKeep, 0755))
	require.NoError(t, os.WriteFile(fileToKeep, []byte("content"), 0644))
	// Add a file inside one of the directories to be removed
	require.NoError(t, os.WriteFile(filepath.Join(dirToRemove1, "file.txt"), []byte("content"), 0644))

	removeDirsWithSuffix(dirToRemove) // Base path is 'data' in 'tempDir'

	assert.NoDirExists(t, dirToRemove, "'data' should be removed")
	assert.NoDirExists(t, dirToRemove1, "'data1' should be removed")
	assert.NoDirExists(t, dirToRemove10, "'data10' should be removed")
	assert.DirExists(t, dirToKeep, "'dataabc' should NOT be removed")
	assert.FileExists(t, fileToKeep, "'data_file.txt' should NOT be removed")
}

// TestStandardWriter_MatchAndFilterOutput_WithDSL_280 tests DSL based matching/filtering.

func TestEvalDslExpr_ErrorScenarios_290(t *testing.T) {
	// Scenario 1: resultToMap fails (hard to simulate directly without corrupting Result struct definition or mapstructure itself)
	// This would be an internal mapstructure error usually.

	// Scenario 2: dsl.EvalExpr returns an error that is NOT ignored by ignoreErr
	result := &Result{Response: &navigation.Response{StatusCode: 200}}
	// This expression is syntactically invalid for dsl
	evalErr := evalDslExpr(result, "status_code ==")
	assert.False(t, evalErr, "Should return false on DSL evaluation error (not ignored)")

	// Scenario 3: dsl.EvalExpr returns an error that IS ignored by ignoreErr
	originalShowDSLErr := showDSLErr
	showDSLErr = false // Ensure errors can be ignored
	defer func() { showDSLErr = originalShowDSLErr }()
	// DSL expression that might cause "No parameter" type error if a helper function was used incorrectly
	// e.g. "contains(header, 'foo')" if 'header' is not defined. 'dsl.EvalExpr' would error.
	// For a direct "No parameter" string error, we'd need dsl.EvalExpr to return that specific string.
	// Let's assume dsl.EvalExpr returns an error whose string contains "No parameter".
	// This requires more intricate knowledge of dsl internals or mocking dsl.EvalExpr.
	// Instead, we test ignoreErr directly.
	// The evalDslExpr uses ignoreErr. If dsl.EvalExpr produces an ignorable error,
	// evalDslExpr should return false (as res would not be true).
	// Example: `dsl.EvalExpr("undefined_var > 10", ...)` might cause an error about `undefined_var`.
	// If this error is ignorable, evalDslExpr returns false.

	evalIgnorableErr := evalDslExpr(result, "non_existent_func(status_code)")
	assert.False(t, evalIgnorableErr, "Should return false if DSL evaluation error is ignored")

	// Scenario 4: dsl.EvalExpr returns false (not an error, just false condition)
	evalFalse := evalDslExpr(result, "status_code == 404")
	assert.False(t, evalFalse, "Should return false if DSL condition is false")
}

// TestStandardWriter_Close_FileErrors_405 aims to test error propagation from
// fileWriter.Close() calls. Making os.File.Close() error reliably is tricky.
// This test simulates a scenario where closing a file descriptor twice might lead to an error,
// which is then caught by fileWriter's Close method.

func TestIgnoreErr_ShowDSLErrTrue_102(t *testing.T) {
	originalShowDSLErr := showDSLErr
	showDSLErr = true // Explicitly set for this test
	defer func() { showDSLErr = originalShowDSLErr }()

	noParamErr := errors.New("something No parameter something")
	assert.False(t, ignoreErr(noParamErr), "Error containing 'No parameter' should NOT be ignored when showDSLErr is true")

	standardErr := errors.New("a standard error")
	assert.False(t, ignoreErr(standardErr), "Standard error should not be ignored when showDSLErr is true")
}

// TestRemoveDirsWithSuffix_270 tests removal of suffixed directories.

