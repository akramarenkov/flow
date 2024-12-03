package inspect

import (
	"testing"

	"github.com/akramarenkov/flow/priority/divider"
	"github.com/stretchr/testify/require"
)

func TestDefaultSet(t *testing.T) {
	require.NotEmpty(t, DefaultSet())
}

func TestDefaultSetFair(t *testing.T) {
	expected := []map[uint]uint{
		{
			1: 100, 2: 100, 3: 100, 4: 100, 5: 100,
			6: 100, 7: 100, 8: 100, 9: 100, 10: 100,
		},
	}

	for id, opts := range DefaultSet() {
		distribution := make(map[uint]uint)

		require.NoError(t, divider.Fair(opts.Quantity, opts.Priorities, distribution))
		require.Equal(t, expected[id], distribution)
	}
}

func TestDefaultSetRate(t *testing.T) {
	expected := []map[uint]uint{
		{
			1: 19, 2: 37, 3: 55, 4: 73, 5: 91, 6: 109, 7: 127, 8: 145,
			9: 163, 10: 181,
		},
	}

	for id, opts := range DefaultSet() {
		distribution := make(map[uint]uint)

		require.NoError(t, divider.Rate(opts.Quantity, opts.Priorities, distribution))
		require.Equal(t, expected[id], distribution)
	}
}
