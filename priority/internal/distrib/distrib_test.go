package distrib

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuantity(t *testing.T) {
	quantity, err := Quantity(nil)
	require.NoError(t, err)
	require.Equal(t, uint(0), quantity)

	quantity, err = Quantity(map[uint]uint{})
	require.NoError(t, err)
	require.Equal(t, uint(0), quantity)

	quantity, err = Quantity(map[uint]uint{3: 3, 2: 5, 1: 11})
	require.NoError(t, err)
	require.Equal(t, uint(19), quantity)

	quantity, err = Quantity(map[uint]uint{3: math.MaxUint, 2: 5, 1: 11})
	require.Error(t, err)
	require.Equal(t, uint(0), quantity)
}

func TestIsFilled(t *testing.T) {
	require.False(t, IsFilled(nil))
	require.False(t, IsFilled(map[uint]uint{}))
	require.False(t, IsFilled(map[uint]uint{3: 0, 2: 0, 1: 0}))
	require.False(t, IsFilled(map[uint]uint{3: 1, 2: 0, 1: 0}))
	require.False(t, IsFilled(map[uint]uint{3: 0, 2: 1, 1: 0}))
	require.False(t, IsFilled(map[uint]uint{3: 0, 2: 0, 1: 1}))
	require.False(t, IsFilled(map[uint]uint{3: 1, 2: 1, 1: 0}))
	require.False(t, IsFilled(map[uint]uint{3: 1, 2: 0, 1: 1}))
	require.False(t, IsFilled(map[uint]uint{3: 0, 2: 1, 1: 1}))
	require.True(t, IsFilled(map[uint]uint{3: 1, 2: 1, 1: 1}))
}

func BenchmarkQuantityReference(b *testing.B) {
	distribution := map[uint]uint{3: 3, 2: 5, 1: 11}

	var quantity uint

	for range b.N {
		quantity = 0

		for _, amount := range distribution {
			quantity += amount
		}
	}

	require.Equal(b, uint(19), quantity)
}

func BenchmarkQuantity(b *testing.B) {
	distribution := map[uint]uint{3: 3, 2: 5, 1: 11}

	var (
		quantity uint
		err      error
	)

	for range b.N {
		quantity, err = Quantity(distribution)
	}

	require.NoError(b, err)
	require.Equal(b, uint(19), quantity)
}

func BenchmarkIsFilled(b *testing.B) {
	distribution := map[uint]uint{3: 3, 2: 5, 1: 11}

	var conclusion bool

	for range b.N {
		conclusion = IsFilled(distribution)
	}

	require.True(b, conclusion)
}
