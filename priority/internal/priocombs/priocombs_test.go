package priocombs

import (
	"math/big"
	"slices"
	"testing"

	"github.com/akramarenkov/safe/intspec"
	"github.com/akramarenkov/seq"
	"github.com/stretchr/testify/require"
)

func TestQuantity(t *testing.T) {
	require.Equal(t, uint64(0), Quantity(nil).Uint64())
	require.Equal(t, uint64(0), Quantity([]uint{}).Uint64())
	require.Equal(t, uint64(1), Quantity(seq.Linear[uint](1, 1)).Uint64())
	require.Equal(t, uint64(3), Quantity(seq.Linear[uint](2, 1)).Uint64())
	require.Equal(t, uint64(7), Quantity(seq.Linear[uint](3, 1)).Uint64())
	require.Equal(t, uint64(15), Quantity(seq.Linear[uint](4, 1)).Uint64())
	require.Equal(t, uint64(31), Quantity(seq.Linear[uint](5, 1)).Uint64())
	require.Equal(t, uint64(63), Quantity(seq.Linear[uint](6, 1)).Uint64())
	require.Equal(t, uint64(127), Quantity(seq.Linear[uint](7, 1)).Uint64())
	require.Equal(t, uint64(255), Quantity(seq.Linear[uint](8, 1)).Uint64())
	require.Equal(t, uint64(511), Quantity(seq.Linear[uint](9, 1)).Uint64())
}

func TestSize(t *testing.T) {
	require.Equal(t, uint64(0), Size(nil))
	require.Equal(t, uint64(0), Size([]uint{}))
	require.Equal(t, uint64(1), Size(seq.Linear[uint](1, 1)))
	require.Equal(t, uint64(3), Size(seq.Linear[uint](2, 1)))
	require.Equal(t, uint64(7), Size(seq.Linear[uint](3, 1)))
	require.Equal(t, uint64(15), Size(seq.Linear[uint](4, 1)))
	require.Equal(t, uint64(31), Size(seq.Linear[uint](5, 1)))
	require.Equal(t, uint64(63), Size(seq.Linear[uint](6, 1)))
	require.Equal(t, uint64(127), Size(seq.Linear[uint](7, 1)))
	require.Equal(t, uint64(255), Size(seq.Linear[uint](8, 1)))
	require.Equal(t, uint64(511), Size(seq.Linear[uint](9, 1)))

	require.Equal(t, uint64(1<<63-1), Size(seq.Linear[uint](63, 1)))
	require.Equal(t, uint64(intspec.MaxUint64), Size(seq.Linear[uint](64, 1)))
}

func TestIterZero(t *testing.T) {
	expected := [][]uint{}

	testIter(t, nil, expected)
	testIter(t, []uint{}, expected)
}

func TestIter1(t *testing.T) {
	expected := [][]uint{
		{1},
	}

	testIter(t, seq.Linear[uint](1, 1), expected)
}

func TestIter21(t *testing.T) {
	expected := [][]uint{
		{2},
		{1},
		{2, 1},
	}

	testIter(t, seq.Linear[uint](2, 1), expected)
}

func TestIter321(t *testing.T) {
	expected := [][]uint{
		{3},
		{2},
		{1},
		{2, 1},
		{3, 2},
		{3, 1},
		{3, 2, 1},
	}

	testIter(t, seq.Linear[uint](3, 1), expected)
}

func TestIter4321(t *testing.T) {
	expected := [][]uint{
		{4},
		{3},
		{2},
		{1},
		{2, 1},
		{3, 2},
		{3, 1},
		{4, 3},
		{4, 2},
		{4, 1},
		{3, 2, 1},
		{4, 3, 1},
		{4, 2, 1},
		{4, 3, 2},
		{4, 3, 2, 1},
	}

	testIter(t, seq.Linear[uint](4, 1), expected)
}

func TestIter54321(t *testing.T) {
	expected := [][]uint{
		{5},
		{4},
		{3},
		{2},
		{1},
		{2, 1},
		{3, 2},
		{3, 1},
		{4, 3},
		{4, 2},
		{4, 1},
		{5, 4},
		{5, 3},
		{5, 2},
		{5, 1},
		{3, 2, 1},
		{4, 3, 1},
		{4, 2, 1},
		{4, 3, 2},
		{4, 3, 2, 1},
		{5, 2, 1},
		{5, 3, 2},
		{5, 3, 1},
		{5, 4, 3},
		{5, 4, 2},
		{5, 4, 1},
		{5, 3, 2, 1},
		{5, 4, 3, 1},
		{5, 4, 2, 1},
		{5, 4, 3, 2},
		{5, 4, 3, 2, 1},
	}

	testIter(t, seq.Linear[uint](5, 1), expected)
}

func TestIter702010(t *testing.T) {
	expected := [][]uint{
		{70},
		{20},
		{10},
		{20, 10},
		{70, 20},
		{70, 10},
		{70, 20, 10},
	}

	testIter(t, []uint{70, 20, 10}, expected)
}

func testIter(t *testing.T, priorities []uint, expected [][]uint) {
	combinations := make([][]uint, 0, Size(priorities))

	for combination := range Iter(priorities) {
		combinations = append(combinations, slices.Clone(combination))
	}

	require.ElementsMatch(t, combinations, expected)
}

func TestIterPartial(t *testing.T) {
	priorities := seq.Linear[uint](15, 1)

	for combination := range Iter(priorities) {
		if len(combination) == 1 {
			return
		}

		require.NotEqual(t, 1, len(combination))
	}
}

func TestIterSize(t *testing.T) {
	priorities := seq.Linear[uint](15, 1)

	size := uint64(0)

	for range Iter(priorities) {
		size++
	}

	require.Equal(t, Size(priorities), size)
}

func TestIterCombinationDecreasing(t *testing.T) {
	priorities := seq.Linear[uint](15, 1)

	for combination := range Iter(priorities) {
		require.IsDecreasing(t, combination, "combination: %v", combination)
	}
}

func BenchmarkQuantity(b *testing.B) {
	priorities := seq.Linear[uint](15, 1)

	b.ResetTimer()

	var quantity *big.Int

	for range b.N {
		quantity = Quantity(priorities)
	}

	require.Equal(b, uint64(1<<15-1), quantity.Uint64())
}

func BenchmarkSize(b *testing.B) {
	priorities := seq.Linear[uint](15, 1)

	b.ResetTimer()

	var size uint64

	for range b.N {
		size = Size(priorities)
	}

	require.Equal(b, uint64(1<<15-1), size)
}

func BenchmarkIter(b *testing.B) {
	priorities := seq.Linear[uint](15, 1)

	b.ResetTimer()

	var combination []uint

	for range b.N {
		for combination = range Iter(priorities) {
			_ = combination
		}
	}

	require.NotNil(b, combination)
}
