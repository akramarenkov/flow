package simple

import (
	"testing"

	"github.com/akramarenkov/flow/priority"
	"github.com/akramarenkov/flow/priority/divider"
	"github.com/akramarenkov/flow/priority/types"

	"github.com/stretchr/testify/require"
)

func TestOptsAddInput(t *testing.T) {
	opts := Opts[int]{
		Divider:          divider.Fair,
		Handle:           func(types.Prioritized[int]) {},
		HandlersQuantity: 6,
	}

	_, err := New(opts)
	require.Error(t, err)

	require.Error(t, opts.AddInput(0, make(chan int)))
	require.Error(t, opts.AddInput(1, nil))
	require.NoError(t, opts.AddInput(1, make(chan int)))
	require.Error(t, opts.AddInput(1, make(chan int)))

	_, err = New(opts)
	require.NoError(t, err)
}

func TestOptsValidation(t *testing.T) {
	opts := Opts[int]{
		Handle: func(types.Prioritized[int]) {},
	}

	_, err := New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Divider: divider.Fair,
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Divider:          divider.Fair,
		Handle:           func(types.Prioritized[int]) {},
		HandlersQuantity: 6,
		Inputs: map[uint]<-chan int{
			1: make(chan int),
		},
	}

	_, err = New(opts)
	require.NoError(t, err)
}

func TestDiscipline(t *testing.T) {
	testDiscipline(t, divider.Fair, false)
}

func TestDisciplineError(t *testing.T) {
	calls := 0

	wrong := func(quantity uint, priorities []uint, distribution map[uint]uint) error {
		calls++

		if calls < 2 {
			return divider.Fair(quantity, priorities, distribution)
		}

		return priority.ErrDividerBad
	}

	testDiscipline(t, wrong, true)
}

func testDiscipline(t *testing.T, divider types.Divider, isErrorExpected bool) {
	handlersQuantity := uint(6)
	itemsQuantity := 100000
	inputCapacity := handlersQuantity

	inputs := map[uint]chan int{
		3: make(chan int, inputCapacity),
		2: make(chan int, inputCapacity),
		1: make(chan int, inputCapacity),
	}

	totalItemsQuantity := itemsQuantity * len(inputs)

	received := make(chan types.Prioritized[int], totalItemsQuantity)

	handle := func(prioritized types.Prioritized[int]) {
		received <- prioritized
	}

	opts := Opts[int]{
		Divider:          divider,
		Handle:           handle,
		HandlersQuantity: handlersQuantity,
	}

	for priority, channel := range inputs {
		err := opts.AddInput(priority, channel)
		require.NoError(t, err)
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	for _, input := range inputs {
		go func() {
			defer close(input)

			for item := range itemsQuantity {
				input <- item
			}
		}()
	}

	err = <-discipline.Err()

	close(received)

	if isErrorExpected {
		require.Error(t, err)
		require.NotEqual(t, totalItemsQuantity, len(received))

		return
	}

	require.NoError(t, err)
	require.Len(t, received, totalItemsQuantity)
}
