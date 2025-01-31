package unmanaged

import (
	"maps"
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiscipline(t *testing.T) {
	testDiscipline(t, map[uint]uint{3: 1e2, 2: 1e2, 1: 1e2}, nil, nil)
}

func TestDisciplineFailAfter(t *testing.T) {
	testDiscipline(
		t,
		map[uint]uint{3: 1e6, 2: 1e6, 1: 1e6},
		map[uint]uint{3: 1e1},
		nil,
	)

	testDiscipline(
		t,
		map[uint]uint{3: 1e6, 2: 1e6, 1: 1e6},
		map[uint]uint{3: 1e1, 2: 1e7, 1: 1e7},
		nil,
	)

	testDiscipline(
		t,
		map[uint]uint{3: 1e6, 2: 1e1, 1: 1e1},
		map[uint]uint{3: 1e6, 2: 1e2, 1: 1e2},
		nil,
	)

	testDiscipline(
		t,
		map[uint]uint{3: 1e2, 2: 1e2, 1: 1e2},
		map[uint]uint{3: 1e3, 2: 1e3, 1: 1e3},
		nil,
	)
}

func TestDisciplineMisses(t *testing.T) {
	testDiscipline(
		t,
		map[uint]uint{3: 1e2, 2: 1e2, 1: 1e2},
		nil,
		map[uint]uint{3: 1, 1: 1e3},
	)

	testDiscipline(
		t,
		map[uint]uint{3: 1e2, 2: 1e2, 1: 1e2},
		map[uint]uint{3: 1},
		map[uint]uint{3: 1},
	)
}

func testDiscipline(
	t *testing.T,
	itemsQuantity map[uint]uint,
	failAfter map[uint]uint,
	misses map[uint]uint,
) {
	inputs := make(map[uint]chan uint, len(itemsQuantity))
	inputsOpts := make(map[uint]<-chan uint, len(itemsQuantity))
	received := make(map[uint]uint, len(itemsQuantity))

	// Input channel capacity is equal to the quantity of data items to keep it simple
	// due to the lack of processing of possible io locks
	for priority, quantity := range itemsQuantity {
		inputs[priority] = make(chan uint, quantity)
		inputsOpts[priority] = inputs[priority]
	}

	// Adds one nil input channel, for coverage
	for priority := range uint(math.MaxUint) {
		if _, exists := inputsOpts[priority]; !exists {
			inputsOpts[priority] = nil
			break
		}
	}

	for priority := range itemsQuantity {
		received[priority] = 0
	}

	opts := Opts[uint]{
		FailAfter:        failAfter,
		HandlersQuantity: 1,
		Inputs:           inputsOpts,
		Misses:           misses,
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	wg := new(sync.WaitGroup)

	for priority, channel := range inputs {
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer close(channel)

			for item := range itemsQuantity[priority] {
				channel <- item
			}
		}()
	}

	for item := range discipline.Output() {
		received[item.Priority]++

		discipline.Release(item.Priority)
	}

	wg.Wait()

	if isErrorExpected(itemsQuantity, failAfter) {
		require.Error(t, <-discipline.Err())
		return
	}

	expected := maps.Clone(itemsQuantity)
	decreaseByMisses(expected, misses)

	require.NoError(t, <-discipline.Err())
	require.Equal(t, expected, received)
}

func isErrorExpected(itemsQuantity, failAfter map[uint]uint) bool {
	for priority, after := range failAfter {
		if after != 0 && after <= itemsQuantity[priority] {
			return true
		}
	}

	return false
}

func decreaseByMisses(target, misses map[uint]uint) {
	for priority := range target {
		if target[priority] < misses[priority] {
			target[priority] = 0
			continue
		}

		target[priority] -= misses[priority]
	}
}

func TestDisciplineError(t *testing.T) {
	opts := Opts[uint]{}

	discipline, err := New(opts)
	require.Error(t, err)
	require.Equal(t, (*Discipline[uint])(nil), discipline)

	opts = Opts[uint]{
		HandlersQuantity: 1,
	}

	discipline, err = New(opts)
	require.Error(t, err)
	require.Equal(t, (*Discipline[uint])(nil), discipline)
}
