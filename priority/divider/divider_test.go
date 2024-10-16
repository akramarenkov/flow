package divider

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFair(t *testing.T) {
	testDivider(t, Fair, []uint{3, 2, 1}, 6, map[uint]uint{3: 2, 2: 2, 1: 2})
	testDivider(t, Fair, []uint{70, 20, 10}, 100, map[uint]uint{70: 34, 20: 33, 10: 33})

	testDivider(t, Fair, []uint{3}, 0, map[uint]uint{3: 0})
	testDivider(t, Fair, []uint{3}, 1, map[uint]uint{3: 1})
	testDivider(t, Fair, []uint{3}, 2, map[uint]uint{3: 2})
	testDivider(t, Fair, []uint{3}, 3, map[uint]uint{3: 3})
	testDivider(t, Fair, []uint{3}, 4, map[uint]uint{3: 4})
	testDivider(t, Fair, []uint{3}, 5, map[uint]uint{3: 5})
	testDivider(t, Fair, []uint{3}, 6, map[uint]uint{3: 6})

	testDivider(t, Fair, []uint{3, 2, 1}, 0, map[uint]uint{3: 0, 2: 0, 1: 0})
	testDivider(t, Fair, []uint{3, 2, 1}, 1, map[uint]uint{3: 1, 2: 0, 1: 0})
	testDivider(t, Fair, []uint{3, 2, 1}, 2, map[uint]uint{3: 1, 2: 1, 1: 0})
	testDivider(t, Fair, []uint{3, 2, 1}, 3, map[uint]uint{3: 1, 2: 1, 1: 1})
	testDivider(t, Fair, []uint{3, 2, 1}, 4, map[uint]uint{3: 2, 2: 1, 1: 1})
	testDivider(t, Fair, []uint{3, 2, 1}, 5, map[uint]uint{3: 2, 2: 2, 1: 1})
	testDivider(t, Fair, []uint{3, 2, 1}, 6, map[uint]uint{3: 2, 2: 2, 1: 2})

	testDivider(t, Fair, []uint{4, 3, 2, 1}, 0, map[uint]uint{4: 0, 3: 0, 2: 0, 1: 0})
	testDivider(t, Fair, []uint{4, 3, 2, 1}, 1, map[uint]uint{4: 1, 3: 0, 2: 0, 1: 0})
	testDivider(t, Fair, []uint{4, 3, 2, 1}, 2, map[uint]uint{4: 1, 3: 1, 2: 0, 1: 0})
	testDivider(t, Fair, []uint{4, 3, 2, 1}, 3, map[uint]uint{4: 1, 3: 1, 2: 1, 1: 0})
	testDivider(t, Fair, []uint{4, 3, 2, 1}, 4, map[uint]uint{4: 1, 3: 1, 2: 1, 1: 1})
	testDivider(t, Fair, []uint{4, 3, 2, 1}, 5, map[uint]uint{4: 2, 3: 1, 2: 1, 1: 1})
	testDivider(t, Fair, []uint{4, 3, 2, 1}, 6, map[uint]uint{4: 2, 3: 2, 2: 1, 1: 1})
	testDivider(t, Fair, []uint{4, 3, 2, 1}, 7, map[uint]uint{4: 2, 3: 2, 2: 2, 1: 1})
	testDivider(t, Fair, []uint{4, 3, 2, 1}, 8, map[uint]uint{4: 2, 3: 2, 2: 2, 1: 2})
}

func TestRate(t *testing.T) {
	testDivider(t, Rate, []uint{3, 2, 1}, 6, map[uint]uint{3: 3, 2: 2, 1: 1})
	testDivider(t, Rate, []uint{70, 20, 10}, 100, map[uint]uint{70: 70, 20: 20, 10: 10})

	testDivider(t, Rate, []uint{3}, 0, map[uint]uint{3: 0})
	testDivider(t, Rate, []uint{3}, 1, map[uint]uint{3: 1})
	testDivider(t, Rate, []uint{3}, 2, map[uint]uint{3: 2})
	testDivider(t, Rate, []uint{3}, 3, map[uint]uint{3: 3})
	testDivider(t, Rate, []uint{3}, 4, map[uint]uint{3: 4})
	testDivider(t, Rate, []uint{3}, 5, map[uint]uint{3: 5})
	testDivider(t, Rate, []uint{3}, 6, map[uint]uint{3: 6})

	testDivider(t, Rate, []uint{3, 2, 1}, 0, map[uint]uint{3: 0, 2: 0, 1: 0})
	testDivider(t, Rate, []uint{3, 2, 1}, 1, map[uint]uint{3: 1, 2: 0, 1: 0})
	testDivider(t, Rate, []uint{3, 2, 1}, 2, map[uint]uint{3: 1, 2: 1, 1: 0})
	testDivider(t, Rate, []uint{3, 2, 1}, 3, map[uint]uint{3: 1, 2: 1, 1: 1})
	testDivider(t, Rate, []uint{3, 2, 1}, 4, map[uint]uint{3: 2, 2: 1, 1: 1})
	testDivider(t, Rate, []uint{3, 2, 1}, 5, map[uint]uint{3: 3, 2: 1, 1: 1})
	testDivider(t, Rate, []uint{3, 2, 1}, 6, map[uint]uint{3: 3, 2: 2, 1: 1})
	testDivider(t, Rate, []uint{3, 2, 1}, 7, map[uint]uint{3: 4, 2: 2, 1: 1})
	testDivider(t, Rate, []uint{3, 2, 1}, 8, map[uint]uint{3: 4, 2: 3, 1: 1})
	testDivider(t, Rate, []uint{3, 2, 1}, 9, map[uint]uint{3: 4, 2: 3, 1: 2})
	testDivider(t, Rate, []uint{3, 2, 1}, 10, map[uint]uint{3: 5, 2: 3, 1: 2})
	testDivider(t, Rate, []uint{3, 2, 1}, 11, map[uint]uint{3: 6, 2: 3, 1: 2})
	testDivider(t, Rate, []uint{3, 2, 1}, 12, map[uint]uint{3: 6, 2: 4, 1: 2})

	testDivider(t, Rate, []uint{4, 3, 2, 1}, 0, map[uint]uint{4: 0, 3: 0, 2: 0, 1: 0})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 1, map[uint]uint{4: 1, 3: 0, 2: 0, 1: 0})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 2, map[uint]uint{4: 1, 3: 1, 2: 0, 1: 0})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 3, map[uint]uint{4: 1, 3: 1, 2: 1, 1: 0})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 4, map[uint]uint{4: 1, 3: 1, 2: 1, 1: 1})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 5, map[uint]uint{4: 2, 3: 1, 2: 1, 1: 1})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 6, map[uint]uint{4: 3, 3: 1, 2: 1, 1: 1})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 7, map[uint]uint{4: 4, 3: 1, 2: 1, 1: 1})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 8, map[uint]uint{4: 4, 3: 2, 2: 1, 1: 1})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 9, map[uint]uint{4: 4, 3: 3, 2: 1, 1: 1})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 10, map[uint]uint{4: 4, 3: 3, 2: 2, 1: 1})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 11, map[uint]uint{4: 5, 3: 3, 2: 2, 1: 1})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 12, map[uint]uint{4: 5, 3: 4, 2: 2, 1: 1})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 13, map[uint]uint{4: 5, 3: 4, 2: 3, 1: 1})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 14, map[uint]uint{4: 5, 3: 4, 2: 3, 1: 2})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 15, map[uint]uint{4: 6, 3: 4, 2: 3, 1: 2})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 16, map[uint]uint{4: 7, 3: 4, 2: 3, 1: 2})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 17, map[uint]uint{4: 8, 3: 4, 2: 3, 1: 2})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 18, map[uint]uint{4: 8, 3: 5, 2: 3, 1: 2})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 19, map[uint]uint{4: 8, 3: 6, 2: 3, 1: 2})
	testDivider(t, Rate, []uint{4, 3, 2, 1}, 20, map[uint]uint{4: 8, 3: 6, 2: 4, 1: 2})
}

func TestRatePanic(t *testing.T) {
	require.Panics(t, func() { testDivider(t, Rate, []uint{math.MaxUint, 1}, 0, nil) })
}

func testDivider(
	t *testing.T,
	divider Divider,
	priorities []uint,
	quantity uint,
	expected map[uint]uint,
) {
	distribution := make(map[uint]uint)

	divider(priorities, quantity, distribution)
	require.Equal(t, expected, distribution)
}

func BenchmarkFair(b *testing.B) {
	priorities := []uint{3, 2, 1}
	distribution := make(map[uint]uint)

	for range b.N {
		// Worst case is when the remainder of dividing the quantity of data handlers
		// by the number of priorities is equal to the maximum value
		Fair(priorities, 8, distribution)
	}

	require.NotNil(b, distribution)
}

func BenchmarkRate(b *testing.B) {
	priorities := []uint{3, 2, 1}
	distribution := make(map[uint]uint)

	for range b.N {
		// Worst case is when the remainder of dividing the quantity of data handlers
		// by the sum of priorities is equal to the maximum value
		Rate(priorities, 11, distribution)
	}

	require.NotNil(b, distribution)
}
