package testutils

import (
	"testing"

	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTestCaseCompareFunc_TargetFound_371 ensures that the CompareFunc
// returns no error when the target string is present in the input 'got' slice.
func TestTestCaseCompareFunc_TargetFound_371(t *testing.T) {
	require.NotEmpty(t, TestCases, "TestCases should not be empty")
	tc := TestCases[0] // Using the first test case defined in the source
	require.NotNil(t, tc.CompareFunc, "CompareFunc should be defined for the test case")

	target := tc.Target
	gotWithTarget := []string{"some other string", "prefix" + target + "suffix", "another string"}

	err := tc.CompareFunc(target, gotWithTarget)
	assert.NoError(t, err, "Expected no error when target is found in 'got' slice")
}

// TestTestCaseCompareFunc_TargetNotFound_824 ensures that the CompareFunc
// returns an error when the target string is not present in the input 'got' slice.
// It also checks the content of the error message.

func TestTestCaseCompareFunc_TargetNotFound_824(t *testing.T) {
	require.NotEmpty(t, TestCases, "TestCases should not be empty")
	tc := TestCases[0] // Using the first test case defined in the source
	require.NotNil(t, tc.CompareFunc, "CompareFunc should be defined for the test case")

	target := tc.Target
	gotWithoutTarget := []string{"string one", "string two", "string three"}

	err := tc.CompareFunc(target, gotWithoutTarget)
	require.Error(t, err, "Expected an error when target is not found in 'got' slice")
	assert.Contains(t, err.Error(), "expected "+target+" target in output", "Error message should contain the target")
	assert.Contains(t, err.Error(), strings.Join(gotWithoutTarget, "\n"), "Error message should contain the 'got' strings")
}

