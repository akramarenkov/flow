package utils

import (
	"testing"

	"github.com/akramarenkov/flow/priority/divider"
	"github.com/akramarenkov/safe"
	"github.com/akramarenkov/seq"

	"github.com/stretchr/testify/require"
)

func TestIsNonFatalConfig(t *testing.T) {
	for quantity := range safe.Iter[uint](0, 2) {
		require.False(
			t,
			IsNonFatalConfig([]uint{3, 2, 1}, divider.Fair, quantity),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](3, 100) {
		require.True(
			t,
			IsNonFatalConfig([]uint{3, 2, 1}, divider.Fair, quantity),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](0, 3) {
		require.False(
			t,
			IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, quantity),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](4, 100) {
		require.True(
			t,
			IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, quantity),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](0, 2) {
		require.False(
			t,
			IsNonFatalConfig([]uint{3, 2, 1}, divider.Rate, quantity),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](3, 100) {
		require.True(
			t,
			IsNonFatalConfig([]uint{3, 2, 1}, divider.Rate, quantity),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](0, 3) {
		require.False(
			t,
			IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Rate, quantity),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](4, 100) {
		require.True(
			t,
			IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Rate, quantity),
			"quantity: %v",
			quantity,
		)
	}
}

func TestIsNonFatalConfigPanic(t *testing.T) {
	require.Panics(
		t,
		func() { IsNonFatalConfig(nil, divider.Fair, 0) },
	)
	require.Panics(
		t,
		func() { IsNonFatalConfig(seq.Int[uint](0, 3), divider.Fair, 0) },
	)
}

func TestPickUpMinNonFatalQuantity(t *testing.T) {
	require.Equal(
		t,
		uint(0),
		PickUpMinNonFatalQuantity([]uint{3, 2, 1}, divider.Fair, 2),
	)

	require.Equal(
		t,
		uint(3),
		PickUpMinNonFatalQuantity([]uint{3, 2, 1}, divider.Fair, 10),
	)

	require.Equal(
		t,
		uint(0),
		PickUpMinNonFatalQuantity([]uint{3, 2, 1}, divider.Rate, 2),
	)

	require.Equal(
		t,
		uint(3),
		PickUpMinNonFatalQuantity([]uint{3, 2, 1}, divider.Rate, 10),
	)
}

func TestPickUpMinNonFatalQuantityPanic(t *testing.T) {
	require.Panics(
		t,
		func() { PickUpMinNonFatalQuantity(nil, divider.Fair, 0) },
	)
	require.Panics(
		t,
		func() { PickUpMinNonFatalQuantity(seq.Int[uint](0, 3), divider.Fair, 0) },
	)
}

func TestPickUpMaxNonFatalQuantity(t *testing.T) {
	require.Equal(
		t,
		uint(0),
		PickUpMaxNonFatalQuantity([]uint{3, 2, 1}, divider.Fair, 2),
	)

	require.Equal(
		t,
		uint(10),
		PickUpMaxNonFatalQuantity([]uint{3, 2, 1}, divider.Fair, 10),
	)

	require.Equal(
		t,
		uint(0),
		PickUpMaxNonFatalQuantity([]uint{3, 2, 1}, divider.Rate, 2),
	)

	require.Equal(
		t,
		uint(10),
		PickUpMaxNonFatalQuantity([]uint{3, 2, 1}, divider.Rate, 10),
	)
}

func TestPickUpMaxNonFatalQuantityPanic(t *testing.T) {
	require.Panics(
		t,
		func() { PickUpMaxNonFatalQuantity(nil, divider.Fair, 0) },
	)
	require.Panics(
		t,
		func() { PickUpMaxNonFatalQuantity(seq.Int[uint](0, 3), divider.Fair, 0) },
	)
}

func TestIsSuitableDiff(t *testing.T) {
	require.True(
		t,
		isSuitableDiff(
			map[uint]uint{3: 300, 2: 200, 1: 100},
			map[uint]uint{3: 3, 2: 2, 1: 1},
			100,
			10,
		),
	)

	require.False(
		t,
		isSuitableDiff(
			map[uint]uint{3: 300, 2: 200, 1: 100},
			map[uint]uint{3: 2, 2: 2, 1: 2},
			100,
			10,
		),
	)
}

func TestIsSuitableConfig(t *testing.T) {
	for quantity := range safe.Iter[uint](0, 5) {
		require.False(
			t,
			IsSuitableConfig([]uint{3, 2, 1}, divider.Fair, quantity, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](20, 100) {
		require.True(
			t,
			IsSuitableConfig([]uint{3, 2, 1}, divider.Fair, quantity, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](0, 11) {
		require.False(
			t,
			IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, quantity, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](30, 100) {
		require.True(
			t,
			IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, quantity, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](0, 11) {
		require.False(
			t,
			IsSuitableConfig([]uint{3, 2, 1}, divider.Rate, quantity, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](28, 100) {
		require.True(
			t,
			IsSuitableConfig([]uint{3, 2, 1}, divider.Rate, quantity, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](0, 31) {
		require.False(
			t,
			IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Rate, quantity, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](55, 100) {
		require.True(
			t,
			IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Rate, quantity, 10),
			"quantity: %v",
			quantity,
		)
	}
}

func TestPickUpMinSuitableQuantity(t *testing.T) {
	require.Equal(
		t,
		uint(0),
		PickUpMinSuitableQuantity([]uint{3, 2, 1}, divider.Fair, 5, 10),
	)

	require.Equal(
		t,
		uint(6),
		PickUpMinSuitableQuantity([]uint{3, 2, 1}, divider.Fair, 10, 10),
	)

	require.Equal(
		t,
		uint(0),
		PickUpMinSuitableQuantity([]uint{3, 2, 1}, divider.Rate, 11, 10),
	)

	require.Equal(
		t,
		uint(12),
		PickUpMinSuitableQuantity([]uint{3, 2, 1}, divider.Rate, 28, 10),
	)
}

func TestPickUpMaxSuitableQuantity(t *testing.T) {
	require.Equal(
		t,
		uint(0),
		PickUpMaxSuitableQuantity([]uint{3, 2, 1}, divider.Fair, 5, 10),
	)

	require.Equal(
		t,
		uint(6),
		PickUpMaxSuitableQuantity([]uint{3, 2, 1}, divider.Fair, 10, 10),
	)

	require.Equal(
		t,
		uint(0),
		PickUpMaxSuitableQuantity([]uint{3, 2, 1}, divider.Rate, 11, 10),
	)

	require.Equal(
		t,
		uint(28),
		PickUpMaxSuitableQuantity([]uint{3, 2, 1}, divider.Rate, 28, 10),
	)
}
