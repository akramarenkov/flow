package priority

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
