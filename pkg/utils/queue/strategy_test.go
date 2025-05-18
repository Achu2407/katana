package queue

import (
	"testing"
)

// TestStrategyString_ExistingStrategy_001 tests the String method for Strategy when the strategy exists in the strategiesMap.
func TestStrategyString_ExistingStrategy_001(t *testing.T) {
	// Arrange
	strategy := BreadthFirst

	// Act
	result := strategy.String()

	// Assert
	expected := "breadth-first"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

// TestStrategyString_NonExistingStrategy_002 tests the String method for Strategy when the strategy does not exist in the strategiesMap.

func TestStrategyString_NonExistingStrategy_002(t *testing.T) {
	// Arrange
	strategy := Strategy(999) // Invalid strategy

	// Act
	result := strategy.String()

	// Assert
	expected := ""
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

