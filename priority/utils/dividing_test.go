package utils

import (
	"testing"

	"github.com/akramarenkov/flow/priority/divider"
	"github.com/stretchr/testify/require"
)

func TestIsCorrectTotalQuantity(t *testing.T) {
	require.True(
		t,
		IsCorrectTotalQuantity(
			[]uint{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			divider.Fair,
			1000,
		),
	)

	require.True(
		t,
		IsCorrectTotalQuantity(
			[]uint{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			divider.Rate,
			1000,
		),
	)
}

func TestIsMonotonic(t *testing.T) {
	require.True(
		t,
		IsMonotonic(
			[]uint{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			divider.Fair,
			1000,
		),
	)

	require.True(
		t,
		IsMonotonic(
			[]uint{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			divider.Rate,
			1000,
		),
	)
}
