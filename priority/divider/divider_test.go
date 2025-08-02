package divider

import (
	"math"
	"testing"

	"github.com/akramarenkov/flow/priority/priodefs"

	"github.com/stretchr/testify/require"
)

func TestFair(t *testing.T) {
	testDivider(t, Fair, 6, []uint{3, 2, 1}, map[uint]uint{3: 2, 2: 2, 1: 2})
	testDivider(t, Fair, 10, []uint{7, 2, 1}, map[uint]uint{7: 4, 2: 3, 1: 3})
	testDivider(t, Fair, 100, []uint{70, 20, 10}, map[uint]uint{70: 34, 20: 33, 10: 33})

	testDivider(t, Fair, 0, []uint{3}, map[uint]uint{3: 0})
	testDivider(t, Fair, 1, []uint{3}, map[uint]uint{3: 1})
	testDivider(t, Fair, 2, []uint{3}, map[uint]uint{3: 2})
	testDivider(t, Fair, 3, []uint{3}, map[uint]uint{3: 3})
	testDivider(t, Fair, 4, []uint{3}, map[uint]uint{3: 4})
	testDivider(t, Fair, 5, []uint{3}, map[uint]uint{3: 5})
	testDivider(t, Fair, 6, []uint{3}, map[uint]uint{3: 6})

	testDivider(t, Fair, 0, []uint{3, 2, 1}, map[uint]uint{3: 0, 2: 0, 1: 0})
	testDivider(t, Fair, 1, []uint{3, 2, 1}, map[uint]uint{3: 1, 2: 0, 1: 0})
	testDivider(t, Fair, 2, []uint{3, 2, 1}, map[uint]uint{3: 1, 2: 1, 1: 0})
	testDivider(t, Fair, 3, []uint{3, 2, 1}, map[uint]uint{3: 1, 2: 1, 1: 1})
	testDivider(t, Fair, 4, []uint{3, 2, 1}, map[uint]uint{3: 2, 2: 1, 1: 1})
	testDivider(t, Fair, 5, []uint{3, 2, 1}, map[uint]uint{3: 2, 2: 2, 1: 1})
	testDivider(t, Fair, 6, []uint{3, 2, 1}, map[uint]uint{3: 2, 2: 2, 1: 2})

	testDivider(t, Fair, 0, []uint{4, 3, 2, 1}, map[uint]uint{4: 0, 3: 0, 2: 0, 1: 0})
	testDivider(t, Fair, 1, []uint{4, 3, 2, 1}, map[uint]uint{4: 1, 3: 0, 2: 0, 1: 0})
	testDivider(t, Fair, 2, []uint{4, 3, 2, 1}, map[uint]uint{4: 1, 3: 1, 2: 0, 1: 0})
	testDivider(t, Fair, 3, []uint{4, 3, 2, 1}, map[uint]uint{4: 1, 3: 1, 2: 1, 1: 0})
	testDivider(t, Fair, 4, []uint{4, 3, 2, 1}, map[uint]uint{4: 1, 3: 1, 2: 1, 1: 1})
	testDivider(t, Fair, 5, []uint{4, 3, 2, 1}, map[uint]uint{4: 2, 3: 1, 2: 1, 1: 1})
	testDivider(t, Fair, 6, []uint{4, 3, 2, 1}, map[uint]uint{4: 2, 3: 2, 2: 1, 1: 1})
	testDivider(t, Fair, 7, []uint{4, 3, 2, 1}, map[uint]uint{4: 2, 3: 2, 2: 2, 1: 1})
	testDivider(t, Fair, 8, []uint{4, 3, 2, 1}, map[uint]uint{4: 2, 3: 2, 2: 2, 1: 2})
}

func TestRate(t *testing.T) {
	testDivider(t, Rate, 6, []uint{3, 2, 1}, map[uint]uint{3: 3, 2: 2, 1: 1})
	testDivider(t, Rate, 10, []uint{7, 2, 1}, map[uint]uint{7: 7, 2: 2, 1: 1})
	testDivider(t, Rate, 100, []uint{70, 20, 10}, map[uint]uint{70: 70, 20: 20, 10: 10})

	testDivider(t, Rate, 0, []uint{3}, map[uint]uint{3: 0})
	testDivider(t, Rate, 1, []uint{3}, map[uint]uint{3: 1})
	testDivider(t, Rate, 2, []uint{3}, map[uint]uint{3: 2})
	testDivider(t, Rate, 3, []uint{3}, map[uint]uint{3: 3})
	testDivider(t, Rate, 4, []uint{3}, map[uint]uint{3: 4})
	testDivider(t, Rate, 5, []uint{3}, map[uint]uint{3: 5})
	testDivider(t, Rate, 6, []uint{3}, map[uint]uint{3: 6})

	testDivider(t, Rate, 0, []uint{3, 2, 1}, map[uint]uint{3: 0, 2: 0, 1: 0})
	testDivider(t, Rate, 1, []uint{3, 2, 1}, map[uint]uint{3: 1, 2: 0, 1: 0})
	testDivider(t, Rate, 2, []uint{3, 2, 1}, map[uint]uint{3: 1, 2: 1, 1: 0})
	testDivider(t, Rate, 3, []uint{3, 2, 1}, map[uint]uint{3: 1, 2: 1, 1: 1})
	testDivider(t, Rate, 4, []uint{3, 2, 1}, map[uint]uint{3: 2, 2: 1, 1: 1})
	testDivider(t, Rate, 5, []uint{3, 2, 1}, map[uint]uint{3: 3, 2: 1, 1: 1})
	testDivider(t, Rate, 6, []uint{3, 2, 1}, map[uint]uint{3: 3, 2: 2, 1: 1})
	testDivider(t, Rate, 7, []uint{3, 2, 1}, map[uint]uint{3: 4, 2: 2, 1: 1})
	testDivider(t, Rate, 8, []uint{3, 2, 1}, map[uint]uint{3: 4, 2: 3, 1: 1})
	testDivider(t, Rate, 9, []uint{3, 2, 1}, map[uint]uint{3: 4, 2: 3, 1: 2})
	testDivider(t, Rate, 10, []uint{3, 2, 1}, map[uint]uint{3: 5, 2: 3, 1: 2})
	testDivider(t, Rate, 11, []uint{3, 2, 1}, map[uint]uint{3: 6, 2: 3, 1: 2})
	testDivider(t, Rate, 12, []uint{3, 2, 1}, map[uint]uint{3: 6, 2: 4, 1: 2})

	testDivider(t, Rate, 0, []uint{4, 3, 2, 1}, map[uint]uint{4: 0, 3: 0, 2: 0, 1: 0})
	testDivider(t, Rate, 1, []uint{4, 3, 2, 1}, map[uint]uint{4: 1, 3: 0, 2: 0, 1: 0})
	testDivider(t, Rate, 2, []uint{4, 3, 2, 1}, map[uint]uint{4: 1, 3: 1, 2: 0, 1: 0})
	testDivider(t, Rate, 3, []uint{4, 3, 2, 1}, map[uint]uint{4: 1, 3: 1, 2: 1, 1: 0})
	testDivider(t, Rate, 4, []uint{4, 3, 2, 1}, map[uint]uint{4: 1, 3: 1, 2: 1, 1: 1})
	testDivider(t, Rate, 5, []uint{4, 3, 2, 1}, map[uint]uint{4: 2, 3: 1, 2: 1, 1: 1})
	testDivider(t, Rate, 6, []uint{4, 3, 2, 1}, map[uint]uint{4: 3, 3: 1, 2: 1, 1: 1})
	testDivider(t, Rate, 7, []uint{4, 3, 2, 1}, map[uint]uint{4: 4, 3: 1, 2: 1, 1: 1})
	testDivider(t, Rate, 8, []uint{4, 3, 2, 1}, map[uint]uint{4: 4, 3: 2, 2: 1, 1: 1})
	testDivider(t, Rate, 9, []uint{4, 3, 2, 1}, map[uint]uint{4: 4, 3: 3, 2: 1, 1: 1})
	testDivider(t, Rate, 10, []uint{4, 3, 2, 1}, map[uint]uint{4: 4, 3: 3, 2: 2, 1: 1})
	testDivider(t, Rate, 11, []uint{4, 3, 2, 1}, map[uint]uint{4: 5, 3: 3, 2: 2, 1: 1})
	testDivider(t, Rate, 12, []uint{4, 3, 2, 1}, map[uint]uint{4: 5, 3: 4, 2: 2, 1: 1})
	testDivider(t, Rate, 13, []uint{4, 3, 2, 1}, map[uint]uint{4: 5, 3: 4, 2: 3, 1: 1})
	testDivider(t, Rate, 14, []uint{4, 3, 2, 1}, map[uint]uint{4: 5, 3: 4, 2: 3, 1: 2})
	testDivider(t, Rate, 15, []uint{4, 3, 2, 1}, map[uint]uint{4: 6, 3: 4, 2: 3, 1: 2})
	testDivider(t, Rate, 16, []uint{4, 3, 2, 1}, map[uint]uint{4: 7, 3: 4, 2: 3, 1: 2})
	testDivider(t, Rate, 17, []uint{4, 3, 2, 1}, map[uint]uint{4: 8, 3: 4, 2: 3, 1: 2})
	testDivider(t, Rate, 18, []uint{4, 3, 2, 1}, map[uint]uint{4: 8, 3: 5, 2: 3, 1: 2})
	testDivider(t, Rate, 19, []uint{4, 3, 2, 1}, map[uint]uint{4: 8, 3: 6, 2: 3, 1: 2})
	testDivider(t, Rate, 20, []uint{4, 3, 2, 1}, map[uint]uint{4: 8, 3: 6, 2: 4, 1: 2})
}

func TestRateError(t *testing.T) {
	require.Error(t, Rate(1, []uint{math.MaxUint, 1}, map[uint]uint{}))
}

func testDivider(
	t *testing.T,
	divider priodefs.Divider,
	quantity uint,
	priorities []uint,
	expected map[uint]uint,
) {
	distribution := make(map[uint]uint)

	require.NoError(t, divider(quantity, priorities, distribution))
	require.Equal(t, expected, distribution)
}

func BenchmarkFair(b *testing.B) {
	priorities := []uint{3, 2, 1}
	distribution := make(map[uint]uint)

	var err error

	for range b.N {
		// Worst case is when the remainder of dividing the quantity of data handlers
		// by the number of priorities is equal to the maximum value
		err = Fair(8, priorities, distribution)
	}

	require.NoError(b, err)
	require.NotNil(b, distribution)
}

func BenchmarkRate(b *testing.B) {
	priorities := []uint{3, 2, 1}
	distribution := make(map[uint]uint)

	var err error

	for range b.N {
		// Worst case is when the remainder of dividing the quantity of data handlers
		// by the sum of priorities is equal to the maximum value
		err = Rate(11, priorities, distribution)
	}

	require.NoError(b, err)
	require.NotNil(b, distribution)
}
