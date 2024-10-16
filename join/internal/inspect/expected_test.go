package inspect

import (
	"math"
	"testing"

	"github.com/akramarenkov/safe"
	"github.com/akramarenkov/seq"
	"github.com/stretchr/testify/require"
)

func TestCalcDescriptionZeroes(t *testing.T) {
	require.Equal(t, description{}, calcDescription(0, 0, 0))
	require.Equal(t, description{}, calcDescription(0, 4, 0))
	require.Equal(t, description{}, calcDescription(12, 0, 0))
	require.Equal(t, description{}, calcDescription(12, 4, 0))

	require.Equal(t, description{}, calcDescription(0, 0, 10))
	require.Equal(t, description{}, calcDescription(0, 4, 10))
	require.Equal(t, description{}, calcDescription(12, 0, 10))
}

func TestCalcDescriptionBlockSize1(t *testing.T) {
	for quantity := range safe.Iter[uint](1, 10) {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 10,
				EffectiveQuantity: quantity,
				Joins:             1,
				RemainderQuantity: 0,
				UnusedJoinSize:    0,
			},
			calcDescription(quantity, 1, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](11, 30) {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 10,
				EffectiveQuantity: 10 * (quantity / 10),
				Joins:             uint(math.Ceil(float64(quantity) / 10)),
				RemainderQuantity: quantity % 10,
				UnusedJoinSize:    0,
			},
			calcDescription(quantity, 1, 10),
			"quantity: %v",
			quantity,
		)
	}
}

func TestCalcDescriptionBlockSize3(t *testing.T) {
	for quantity := range safe.Iter[uint](1, 9) {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 9,
				EffectiveQuantity: quantity,
				Joins:             1,
				RemainderQuantity: 0,
				UnusedJoinSize:    1,
			},
			calcDescription(quantity, 3, 10),
			"quantity: %v",
			quantity,
		)
	}

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 9,
			EffectiveQuantity: 9,
			Joins:             1,
			RemainderQuantity: 1,
			UnusedJoinSize:    1,
		},
		calcDescription(10, 3, 10),
	)

	for quantity := range safe.Iter[uint](11, 18) {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 9,
				EffectiveQuantity: 9 * (quantity / 9),
				Joins:             uint(math.Ceil(float64(quantity) / 9)),
				RemainderQuantity: quantity % 9,
				UnusedJoinSize:    1,
			},
			calcDescription(quantity, 3, 10),
			"quantity: %v",
			quantity,
		)
	}

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 9,
			EffectiveQuantity: 18,
			Joins:             2,
			RemainderQuantity: 1,
			UnusedJoinSize:    1,
		},
		calcDescription(19, 3, 10),
	)

	for quantity := range safe.Iter[uint](20, 27) {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 9,
				EffectiveQuantity: 9 * (quantity / 9),
				Joins:             uint(math.Ceil(float64(quantity) / 9)),
				RemainderQuantity: quantity % 9,
				UnusedJoinSize:    1,
			},
			calcDescription(quantity, 3, 10),
			"quantity: %v",
			quantity,
		)
	}

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 9,
			EffectiveQuantity: 27,
			Joins:             3,
			RemainderQuantity: 1,
			UnusedJoinSize:    1,
		},
		calcDescription(28, 3, 10),
	)
}

func TestCalcDescriptionBlockSize4(t *testing.T) {
	for quantity := range safe.Iter[uint](1, 8) {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 8,
				EffectiveQuantity: quantity,
				Joins:             1,
				RemainderQuantity: 0,
				UnusedJoinSize:    2,
			},
			calcDescription(quantity, 4, 10),
			"quantity: %v",
			quantity,
		)
	}

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 8,
			EffectiveQuantity: 8,
			Joins:             1,
			RemainderQuantity: 1,
			UnusedJoinSize:    2,
		},
		calcDescription(9, 4, 10),
	)

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 8,
			EffectiveQuantity: 8,
			Joins:             1,
			RemainderQuantity: 2,
			UnusedJoinSize:    2,
		},
		calcDescription(10, 4, 10),
	)

	for quantity := range safe.Iter[uint](11, 16) {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 8,
				EffectiveQuantity: 8 * (quantity / 8),
				Joins:             uint(math.Ceil(float64(quantity) / 8)),
				RemainderQuantity: quantity % 8,
				UnusedJoinSize:    2,
			},
			calcDescription(quantity, 4, 10),
			"quantity: %v",
			quantity,
		)
	}

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 8,
			EffectiveQuantity: 16,
			Joins:             2,
			RemainderQuantity: 1,
			UnusedJoinSize:    2,
		},
		calcDescription(17, 4, 10),
	)

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 8,
			EffectiveQuantity: 16,
			Joins:             2,
			RemainderQuantity: 2,
			UnusedJoinSize:    2,
		},
		calcDescription(18, 4, 10),
	)

	for quantity := range safe.Iter[uint](19, 24) {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 8,
				EffectiveQuantity: 8 * (quantity / 8),
				Joins:             uint(math.Ceil(float64(quantity) / 8)),
				RemainderQuantity: quantity % 8,
				UnusedJoinSize:    2,
			},
			calcDescription(quantity, 4, 10),
			"quantity: %v",
			quantity,
		)
	}

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 8,
			EffectiveQuantity: 24,
			Joins:             3,
			RemainderQuantity: 1,
			UnusedJoinSize:    2,
		},
		calcDescription(25, 4, 10),
	)

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 8,
			EffectiveQuantity: 24,
			Joins:             3,
			RemainderQuantity: 2,
			UnusedJoinSize:    2,
		},
		calcDescription(26, 4, 10),
	)
}

func TestCalcDescriptionBlockSize10(t *testing.T) {
	for quantity := range safe.Iter[uint](1, 10) {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 10,
				EffectiveQuantity: quantity,
				Joins:             1,
				RemainderQuantity: 0,
				UnusedJoinSize:    0,
			},
			calcDescription(quantity, 10, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](11, 30) {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 10,
				EffectiveQuantity: 10 * (quantity / 10),
				Joins:             uint(math.Ceil(float64(quantity) / 10)),
				RemainderQuantity: quantity % 10,
				UnusedJoinSize:    0,
			},
			calcDescription(quantity, 10, 10),
			"quantity: %v",
			quantity,
		)
	}
}

func TestCalcDescriptionBlockSize11(t *testing.T) {
	for quantity := range safe.Iter[uint](1, 11) {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 11,
				EffectiveQuantity: quantity,
				Joins:             1,
				RemainderQuantity: 0,
				UnusedJoinSize:    0,
			},
			calcDescription(quantity, 11, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](12, 40) {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 11,
				EffectiveQuantity: 11 * (quantity / 11),
				Joins:             uint(math.Ceil(float64(quantity) / 11)),
				RemainderQuantity: quantity % 11,
				UnusedJoinSize:    0,
			},
			calcDescription(quantity, 11, 10),
			"quantity: %v",
			quantity,
		)
	}
}

func TestExpectedZeroes(t *testing.T) {
	require.Equal(t, [][]uint{}, Expected(0, 0, 0))
	require.Equal(t, [][]uint{}, Expected(0, 4, 0))
	require.Equal(t, [][]uint{}, Expected(12, 0, 0))
	require.Equal(t, [][]uint{}, Expected(12, 4, 0))

	require.Equal(t, [][]uint{}, Expected(0, 0, 10))
	require.Equal(t, [][]uint{}, Expected(0, 4, 10))
	require.Equal(t, [][]uint{}, Expected(12, 0, 10))
}

func TestExpectedBlockSize1(t *testing.T) {
	for quantity := range safe.Iter[uint](1, 10) {
		require.Equal(
			t,
			[][]uint{
				seq.Int[uint](1, quantity),
			},
			Expected(quantity, 1, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](11, 20) {
		require.Equal(
			t,
			[][]uint{
				seq.Int[uint](1, 10),
				seq.Int[uint](11, quantity),
			},
			Expected(quantity, 1, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](21, 30) {
		require.Equal(
			t,
			[][]uint{
				seq.Int[uint](1, 10),
				seq.Int[uint](11, 20),
				seq.Int[uint](21, quantity),
			},
			Expected(quantity, 1, 10),
			"quantity: %v",
			quantity,
		)
	}
}

func TestExpectedBlockSize3(t *testing.T) {
	for quantity := range safe.Iter[uint](1, 10) {
		require.Equal(
			t,
			[][]uint{
				seq.Int[uint](1, quantity),
			},
			Expected(quantity, 3, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](11, 19) {
		require.Equal(
			t,
			[][]uint{
				seq.Int[uint](1, 9),
				seq.Int[uint](10, quantity),
			},
			Expected(quantity, 3, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](20, 28) {
		require.Equal(
			t,
			[][]uint{
				seq.Int[uint](1, 9),
				seq.Int[uint](10, 18),
				seq.Int[uint](19, quantity),
			},
			Expected(quantity, 3, 10),
			"quantity: %v",
			quantity,
		)
	}
}

func TestExpectedBlockSize4(t *testing.T) {
	for quantity := range safe.Iter[uint](1, 10) {
		require.Equal(
			t,
			[][]uint{
				seq.Int[uint](1, quantity),
			},
			Expected(quantity, 4, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](11, 18) {
		require.Equal(
			t,
			[][]uint{
				seq.Int[uint](1, 8),
				seq.Int[uint](9, quantity),
			},
			Expected(quantity, 4, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](19, 26) {
		require.Equal(
			t,
			[][]uint{
				seq.Int[uint](1, 8),
				seq.Int[uint](9, 16),
				seq.Int[uint](17, quantity),
			},
			Expected(quantity, 4, 10),
			"quantity: %v",
			quantity,
		)
	}
}

func TestExpectedBlockSize10(t *testing.T) {
	for quantity := range safe.Iter[uint](1, 10) {
		require.Equal(
			t,
			[][]uint{
				seq.Int[uint](1, quantity),
			},
			Expected(quantity, 10, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](11, 20) {
		require.Equal(
			t,
			[][]uint{
				seq.Int[uint](1, 10),
				seq.Int[uint](11, quantity),
			},
			Expected(quantity, 10, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](21, 30) {
		require.Equal(
			t,
			[][]uint{
				seq.Int[uint](1, 10),
				seq.Int[uint](11, 20),
				seq.Int[uint](21, quantity),
			},
			Expected(quantity, 10, 10),
			"quantity: %v",
			quantity,
		)
	}
}

func TestExpectedBlockSize11(t *testing.T) {
	for quantity := range safe.Iter[uint](1, 11) {
		require.Equal(
			t,
			[][]uint{
				seq.Int[uint](1, quantity),
			},
			Expected(quantity, 11, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](12, 22) {
		require.Equal(
			t,
			[][]uint{
				seq.Int[uint](1, 11),
				seq.Int[uint](12, quantity),
			},
			Expected(quantity, 11, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := range safe.Iter[uint](23, 33) {
		require.Equal(
			t,
			[][]uint{
				seq.Int[uint](1, 11),
				seq.Int[uint](12, 22),
				seq.Int[uint](23, quantity),
			},
			Expected(quantity, 11, 10),
			"quantity: %v",
			quantity,
		)
	}
}

func BenchmarkExpected(b *testing.B) {
	quantity, err := safe.IToI[uint](b.N)
	require.NoError(b, err)

	for range b.N {
		_ = Expected(quantity, 4, 10)
	}
}
