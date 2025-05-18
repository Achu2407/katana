package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUniqueContent_Unique_321 tests the behavior of UniqueContent when the content is unique.
func TestUniqueContent_Unique_321(t *testing.T) {
	// Arrange
	simple, err := NewSimple()
	require.NoError(t, err)
	content := []byte("unique content")

	// Act
	isUnique := simple.UniqueContent(content)

	// Assert
	assert.True(t, isUnique)
}

// TestUniqueContent_NotUnique_654 tests the behavior of UniqueContent when the content is not unique.

func TestIsCycle_ExceedsMaxChromeURLLength_987(t *testing.T) {
	// Arrange
	simple, err := NewSimple()
	require.NoError(t, err)
	longURL := string(make([]byte, MaxChromeURLLength+1))

	// Act
	isCycle := simple.IsCycle(longURL)

	// Assert
	assert.True(t, isCycle)
}

func TestUniqueContent_NotUnique_654(t *testing.T) {
	// Arrange
	simple, err := NewSimple()
	require.NoError(t, err)
	content := []byte("duplicate content")
	simple.UniqueContent(content) // Add content first

	// Act
	isUnique := simple.UniqueContent(content)

	// Assert
	assert.False(t, isUnique)
}

// TestIsCycle_ExceedsMaxChromeURLLength_987 tests the behavior of IsCycle when the URL exceeds MaxChromeURLLength.

