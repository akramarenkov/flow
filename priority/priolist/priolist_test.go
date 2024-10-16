package priolist

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompare(t *testing.T) {
	priorities := []uint{2, 1, 3, 5, 4, 3}
	expected := []uint{5, 4, 3, 3, 2, 1}

	slices.SortFunc(priorities, Compare)
	require.Equal(t, expected, priorities)
}

func TestIsValid(t *testing.T) {
	require.Error(t, IsValid(nil))
	require.Error(t, IsValid([]uint{}))
	require.Error(t, IsValid([]uint{3, 1, 2}))
	require.Error(t, IsValid([]uint{3, 2, 0}))
	require.NoError(t, IsValid([]uint{3, 2, 1}))
}
