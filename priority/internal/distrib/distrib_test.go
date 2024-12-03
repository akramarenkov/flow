package distrib

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuantity(t *testing.T) {
	quantity, err := Quantity(nil, nil)
	require.NoError(t, err)
	require.Equal(t, uint(0), quantity)

	quantity, err = Quantity(nil, map[uint]uint{})
	require.NoError(t, err)
	require.Equal(t, uint(0), quantity)

	quantity, err = Quantity(nil, map[uint]uint{3: 3, 2: 5, 1: 11})
	require.NoError(t, err)
	require.Equal(t, uint(0), quantity)

	quantity, err = Quantity([]uint{}, nil)
	require.NoError(t, err)
	require.Equal(t, uint(0), quantity)

	quantity, err = Quantity([]uint{}, map[uint]uint{})
	require.NoError(t, err)
	require.Equal(t, uint(0), quantity)

	quantity, err = Quantity([]uint{}, map[uint]uint{3: 3, 2: 5, 1: 11})
	require.NoError(t, err)
	require.Equal(t, uint(0), quantity)

	quantity, err = Quantity([]uint{3, 2, 1}, nil)
	require.NoError(t, err)
	require.Equal(t, uint(0), quantity)

	quantity, err = Quantity([]uint{3, 2, 1}, map[uint]uint{})
	require.NoError(t, err)
	require.Equal(t, uint(0), quantity)

	quantity, err = Quantity([]uint{3, 2, 1}, map[uint]uint{3: 3, 2: 5, 1: 11})
	require.NoError(t, err)
	require.Equal(t, uint(19), quantity)

	quantity, err = Quantity([]uint{3, 2, 1}, map[uint]uint{3: math.MaxUint, 2: 5, 1: 11})
	require.Error(t, err)
	require.Equal(t, uint(0), quantity)
}

func TestIsFilled(t *testing.T) {
	require.False(t, IsFilled(nil, nil))
	require.False(t, IsFilled(nil, map[uint]uint{}))
	require.False(t, IsFilled(nil, map[uint]uint{3: 1, 2: 1, 1: 1}))
	require.False(t, IsFilled([]uint{}, nil))
	require.False(t, IsFilled([]uint{}, map[uint]uint{}))
	require.False(t, IsFilled([]uint{}, map[uint]uint{3: 1, 2: 1, 1: 1}))
	require.False(t, IsFilled([]uint{3, 2, 1}, nil))
	require.False(t, IsFilled([]uint{3, 2, 1}, map[uint]uint{}))
	require.False(t, IsFilled([]uint{3, 2, 1}, map[uint]uint{3: 0, 2: 0, 1: 0}))
	require.False(t, IsFilled([]uint{3, 2, 1}, map[uint]uint{3: 1, 2: 0, 1: 0}))
	require.False(t, IsFilled([]uint{3, 2, 1}, map[uint]uint{3: 0, 2: 1, 1: 0}))
	require.False(t, IsFilled([]uint{3, 2, 1}, map[uint]uint{3: 0, 2: 0, 1: 1}))
	require.False(t, IsFilled([]uint{3, 2, 1}, map[uint]uint{3: 1, 2: 1, 1: 0}))
	require.False(t, IsFilled([]uint{3, 2, 1}, map[uint]uint{3: 1, 2: 0, 1: 1}))
	require.False(t, IsFilled([]uint{3, 2, 1}, map[uint]uint{3: 0, 2: 1, 1: 1}))
	require.True(t, IsFilled([]uint{3, 2, 1}, map[uint]uint{3: 1, 2: 1, 1: 1}))
}

func BenchmarkQuantityReference(b *testing.B) {
	priorities := []uint{3, 2, 1}
	distribution := map[uint]uint{3: 3, 2: 5, 1: 11}

	var quantity uint

	for range b.N {
		quantity = 0

		for _, priority := range priorities {
			quantity += distribution[priority]
		}
	}

	require.Equal(b, uint(19), quantity)
}

func BenchmarkQuantity(b *testing.B) {
	priorities := []uint{3, 2, 1}
	distribution := map[uint]uint{3: 3, 2: 5, 1: 11}

	var (
		quantity uint
		err      error
	)

	for range b.N {
		quantity, err = Quantity(priorities, distribution)
	}

	require.NoError(b, err)
	require.Equal(b, uint(19), quantity)
}

func BenchmarkIsFilled(b *testing.B) {
	priorities := []uint{3, 2, 1}
	distribution := map[uint]uint{3: 3, 2: 5, 1: 11}

	var conclusion bool

	for range b.N {
		conclusion = IsFilled(priorities, distribution)
	}

	require.True(b, conclusion)
}
