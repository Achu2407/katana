package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestProcess_MissingArguments_001 tests the behavior of the `process` function when no command-line arguments are provided.
func TestProcess_ValidInput_002(t *testing.T) {
	// Arrange
	inputFile := "test_input.txt"
	inputHeadlessFile := "test_input_headless.txt"
	os.Args = []string{"crawl-maze-score", inputFile, inputHeadlessFile}

	// Create test input files
	err := os.WriteFile(inputFile, []byte("/html/body/img/src.found\n/html/body/img/srcset1x.found\n"), 0644)
	require.NoError(t, err)
	defer os.Remove(inputFile)

	err = os.WriteFile(inputHeadlessFile, []byte("/html/body/img/src.found\n/html/body/img/srcset2x.found\n"), 0644)
	require.NoError(t, err)
	defer os.Remove(inputHeadlessFile)

	// Act
	err = process()

	// Assert
	require.NoError(t, err)
}

func TestProcess_MissingArguments_001(t *testing.T) {
	// Arrange
	os.Args = []string{"crawl-maze-score"}

	// Act
	err := process()

	// Assert
	require.NoError(t, err)
}

// TestProcess_ValidInput_002 tests the behavior of the `process` function with valid input files.

