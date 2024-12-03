package measuring

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompareItem(t *testing.T) {
	measurements := []Measure{
		{Item: 3},
		{Item: 1},
		{Item: 4},
		{Item: 2},
	}

	expected := []Measure{
		{Item: 1},
		{Item: 2},
		{Item: 3},
		{Item: 4},
	}

	slices.SortFunc(measurements, CompareItem)
	require.Equal(t, expected, measurements)
}

func TestCompareTime(t *testing.T) {
	measurements := []Measure{
		{Time: 3},
		{Time: 1},
		{Time: 4},
		{Time: 2},
	}

	expected := []Measure{
		{Time: 1},
		{Time: 2},
		{Time: 3},
		{Time: 4},
	}

	slices.SortFunc(measurements, CompareTime)
	require.Equal(t, expected, measurements)
}

func TestKeepReceived(t *testing.T) {
	measurements := []Measure{
		{Item: 1, Kind: KindCompleted},
		{Item: 1, Kind: KindProcessed},
		{Item: 1, Kind: KindReceived},
		{Item: 2, Kind: KindProcessed},
		{Item: 2, Kind: KindCompleted},
		{Item: 2, Kind: KindReceived},
		{Item: 3, Kind: KindProcessed},
		{Item: 3, Kind: KindReceived},
		{Item: 3, Kind: KindCompleted},
	}

	expected := []Measure{
		{Item: 1, Kind: KindReceived},
		{Item: 2, Kind: KindReceived},
		{Item: 3, Kind: KindReceived},
	}

	filtered := slices.DeleteFunc(slices.Clone(measurements), KeepReceived)
	require.Equal(t, expected, filtered)
}
