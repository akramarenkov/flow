package priority

import (
	"math"
	"testing"

	"github.com/akramarenkov/flow/priority/divider"

	"github.com/stretchr/testify/require"
)

func TestDivide(t *testing.T) {
	wrong := func(quantity uint, priorities []uint, distribution map[uint]uint) error {
		for _, priority := range priorities {
			distribution[priority] = quantity
		}

		return nil
	}

	require.NoError(t, divide(divider.Rate, 6, []uint{3, 2, 1}, make(map[uint]uint)))
	require.NoError(t, divide(divider.Rate, 0, []uint{3, 2, 1}, make(map[uint]uint)))
	require.Error(t, divide(divider.Rate, 6, []uint{math.MaxUint, 2, 1}, make(map[uint]uint)))
	require.Error(t, divide(wrong, math.MaxUint, []uint{3, 2, 1}, make(map[uint]uint)))
	require.Error(t, divide(wrong, 6, []uint{3, 2, 1}, make(map[uint]uint)))
}
