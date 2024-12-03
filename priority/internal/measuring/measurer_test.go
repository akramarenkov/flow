package measuring

import (
	"testing"
	"time"

	"github.com/akramarenkov/flow/priority/internal/unmanaged"
	"github.com/stretchr/testify/require"
)

func TestMeasurer(t *testing.T) {
	msr, err := NewMeasurer(6)
	require.NoError(t, err)

	msr.AddWrite(1, 1000)

	msr.AddWrite(2, 500)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 1*time.Second)
	msr.AddWrite(2, 500)

	msr.AddWrite(3, 500)
	msr.AddWaitDevastation(3)
	msr.AddDelay(3, 1*time.Second)
	msr.AddWrite(3, 500)

	msr.SetProcessingDuration(1, 1*time.Millisecond)
	msr.SetProcessingDuration(2, 1*time.Millisecond)
	msr.SetProcessingDuration(3, 1*time.Millisecond)

	opts := unmanaged.Opts[uint]{
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := unmanaged.New(opts)
	require.NoError(t, err)

	_, err = msr.Play(discipline)
	require.NoError(t, err)
}

func TestMeasurerWriteWithDelay(t *testing.T) {
	msr, err := NewMeasurer(6)
	require.NoError(t, err)

	msr.AddWriteWithDelay(1, 1000, 1*time.Millisecond)
	msr.AddWriteWithDelay(2, 1000, 1*time.Millisecond)
	msr.AddWriteWithDelay(3, 1000, 1*time.Millisecond)

	opts := unmanaged.Opts[uint]{
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := unmanaged.New(opts)
	require.NoError(t, err)

	_, err = msr.Play(discipline)
	require.NoError(t, err)
}

func TestMeasurerBufferedInput(t *testing.T) {
	msr, err := NewMeasurer(6)
	require.NoError(t, err)

	msr.AddWrite(1, 1000)
	msr.AddWrite(2, 1000)
	msr.AddWrite(3, 1000)

	inputs := msr.Inputs()

	require.Len(t, inputs, 3)

	for _, channel := range inputs {
		require.Equal(t, msr.HandlersQuantity(), uint(cap(channel)))
	}
}

func TestMeasurerUnbufferedInput(t *testing.T) {
	msr, err := NewMeasurer(6, 0)
	require.NoError(t, err)

	msr.AddWrite(1, 1000)
	msr.AddWrite(2, 1000)
	msr.AddWrite(3, 1000)

	inputs := msr.Inputs()

	require.Len(t, inputs, 3)

	for _, channel := range inputs {
		require.Equal(t, 0, cap(channel))
	}
}

func TestMeasurerPartially(t *testing.T) {
	msr, err := NewMeasurer(6)
	require.NoError(t, err)

	msr.AddWrite(1, 1000)
	msr.AddWrite(2, 1000)
	msr.AddWrite(3, 0)

	opts := unmanaged.Opts[uint]{
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := unmanaged.New(opts)
	require.NoError(t, err)

	_, err = msr.Play(discipline)
	require.NoError(t, err)
}

func TestMeasurerFail(t *testing.T) {
	msr, err := NewMeasurer(6)
	require.NoError(t, err)

	msr.AddWrite(1, 1000)
	msr.AddWrite(2, 1000)
	msr.AddWrite(3, 1000)

	opts := unmanaged.Opts[uint]{
		FailAfter:        map[uint]uint{3: 1},
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := unmanaged.New(opts)
	require.NoError(t, err)

	_, err = msr.Play(discipline)
	require.Error(t, err)
}

func TestMeasurerFailAtWaitDevastation(t *testing.T) {
	msr, err := NewMeasurer(6)
	require.NoError(t, err)

	msr.AddWrite(1, 501)
	msr.AddWaitDevastation(1)
	msr.AddWrite(1, 500)

	msr.AddWrite(2, 501)
	msr.AddWaitDevastation(2)
	msr.AddWrite(2, 500)

	msr.AddWrite(3, 501)
	msr.AddWaitDevastation(3)
	msr.AddWrite(3, 500)

	opts := unmanaged.Opts[uint]{
		FailAfter:        map[uint]uint{3: 500, 2: 500, 1: 500},
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := unmanaged.New(opts)
	require.NoError(t, err)

	_, err = msr.Play(discipline)
	require.Error(t, err)
}

func TestMeasurerFailByMisses(t *testing.T) {
	msr, err := NewMeasurer(6)
	require.NoError(t, err)

	msr.AddWrite(1, 1000)
	msr.AddWrite(2, 1000)
	msr.AddWrite(3, 1000)

	opts := unmanaged.Opts[uint]{
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
		Misses:           map[uint]uint{3: 1},
	}

	discipline, err := unmanaged.New(opts)
	require.NoError(t, err)

	_, err = msr.Play(discipline)
	require.Error(t, err)
}

func TestMeasurerError(t *testing.T) {
	msr, err := NewMeasurer(0)
	require.Error(t, err)
	require.Equal(t, (*Measurer)(nil), msr)
}
