package measuring

import (
	"testing"

	"github.com/akramarenkov/flow/priority/internal/unmanaged"

	"github.com/stretchr/testify/require"
)

func TestBenchmarker(t *testing.T) {
	testBenchmarker(t)
}

func TestBenchmarkerUnbuffered(t *testing.T) {
	testBenchmarker(t, 0)
}

func testBenchmarker(t *testing.T, inputCapacity ...uint) {
	itemsQuantity := uint(1000)

	bnch, err := NewBenchmarker(100, inputCapacity...)
	require.NoError(t, err)

	bnch.AddItems(3, itemsQuantity)
	bnch.AddItems(2, itemsQuantity)
	bnch.AddItems(1, itemsQuantity)

	opts := unmanaged.Opts[uint]{
		HandlersQuantity: bnch.HandlersQuantity(),
		Inputs:           bnch.Inputs(),
	}

	discipline, err := unmanaged.New(opts)
	require.NoError(t, err)

	bnch.Play(discipline)
}

func TestBenchmarkerInputCapacity(t *testing.T) {
	bnch, err := NewBenchmarker(100)
	require.NoError(t, err)

	bnch.AddItems(3, 1000)
	bnch.AddItems(2, 999)
	bnch.AddItems(1, 998)

	inputs := bnch.Inputs()

	require.Equal(t, 1000, cap(inputs[3]))
	require.Equal(t, 999, cap(inputs[2]))
	require.Equal(t, 998, cap(inputs[1]))
}

func TestBenchmarkerInputCapacityUnbuffered(t *testing.T) {
	bnch, err := NewBenchmarker(100, 0)
	require.NoError(t, err)

	bnch.AddItems(3, 1000)
	bnch.AddItems(2, 999)
	bnch.AddItems(1, 998)

	inputs := bnch.Inputs()

	require.Equal(t, 0, cap(inputs[3]))
	require.Equal(t, 0, cap(inputs[2]))
	require.Equal(t, 0, cap(inputs[1]))
}

func TestBenchmarkerInputCapacityIntermediate(t *testing.T) {
	bnch, err := NewBenchmarker(100, 20)
	require.NoError(t, err)

	bnch.AddItems(3, 1000)
	bnch.AddItems(2, 999)
	bnch.AddItems(1, 998)

	inputs := bnch.Inputs()

	require.Equal(t, 20, cap(inputs[3]))
	require.Equal(t, 20, cap(inputs[2]))
	require.Equal(t, 20, cap(inputs[1]))
}

func TestBenchmarkerError(t *testing.T) {
	bnch, err := NewBenchmarker(0)
	require.Error(t, err)
	require.Equal(t, (*Benchmarker)(nil), bnch)
}
