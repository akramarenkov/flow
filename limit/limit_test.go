package limit

import (
	"testing"
	"time"

	"github.com/akramarenkov/safe"
	"github.com/stretchr/testify/require"
)

func TestOptsValidation(t *testing.T) {
	opts := Opts[int]{}

	_, err := New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input: make(chan int),
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input: make(chan int),
		Limit: Rate{
			Interval: time.Second,
			Quantity: 1,
		},
	}

	_, err = New(opts)
	require.NoError(t, err)
}

func TestDiscipline(t *testing.T) {
	limit := Rate{
		Interval: time.Second,
		Quantity: 1000,
	}

	duration := testDiscipline(t, 10000, limit)
	expected := expectedDuration(t, 10000, limit)
	require.InEpsilon(t, expected, duration, 0.1)

	duration = testDiscipline(t, 9999, limit)
	expected = expectedDuration(t, 9999, limit)
	require.InEpsilon(t, expected, duration, 0.1)
}

func testDiscipline(t *testing.T, quantity int, limit Rate) time.Duration {
	input := make(chan int, quantity)

	opts := Opts[int]{
		Input: input,
		Limit: limit,
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	inpSequence := make([]int, 0, quantity)
	outSequence := make([]int, 0, quantity)

	startedAt := time.Now()

	go func() {
		defer close(input)

		for item := range quantity {
			inpSequence = append(inpSequence, item)

			input <- item
		}
	}()

	for item := range discipline.Output() {
		outSequence = append(outSequence, item)
	}

	duration := time.Since(startedAt)

	require.Equal(t, inpSequence, outSequence)
	require.Len(t, outSequence, quantity)

	return duration
}

func expectedDuration(t *testing.T, quantity int, limit Rate) time.Duration {
	inputQuantity := time.Duration(quantity)

	limitQuantity, err := safe.IToI[time.Duration](limit.Quantity)
	require.NoError(t, err)

	// Accuracy of calculations is deliberately roughened (first division is performed
	// and only then multiplication) because such a calculation corresponds to the work
	// of the discipline when closing the input channel: if the quantity of data items
	// written to the input channel is not a multiple of the Quantity field in rate
	// limit structure, then the delay after the transmission of the last data is not
	// performed
	ratio := inputQuantity / limitQuantity

	return ratio * limit.Interval
}

func BenchmarkDisciplineInputCapacity0(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 0)
}

func BenchmarkDisciplineInputCapacity1e0(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 1e0)
}

func BenchmarkDisciplineInputCapacity1e1(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 1e1)
}

func BenchmarkDisciplineInputCapacity1e2(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 1e2)
}

func BenchmarkDisciplineInputCapacity1e3(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 1e3)
}

func BenchmarkDisciplineInputCapacity1e4(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 1e4)
}

func BenchmarkDisciplineInputCapacity1e5(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 1e5)
}

func BenchmarkDisciplineInputCapacity1e6(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 1e6)
}

func BenchmarkDisciplineInputCapacity1e7(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 1e7)
}

func BenchmarkDisciplineInputCapacityQuantity(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, -1)
}

// This benchmark is used to test the impact of input channel capacity on
// performance. Therefore, the value of Quantity field in rate limit structure is
// always set to be greater than the quantity of data items written to the input
// channel so that there is no delay after data transfer.
func benchmarkDisciplineInputCapacity(b *testing.B, capacity int) {
	quantity := b.N

	limitQuantity, err := safe.IToI[uint64](quantity)
	require.NoError(b, err)

	limit := Rate{
		Interval: time.Second,
		Quantity: limitQuantity + 1,
	}

	if capacity < 0 {
		capacity = quantity
	}

	input := make(chan int, capacity)

	opts := Opts[int]{
		Input: input,
		Limit: limit,
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	b.ResetTimer()

	go func() {
		defer close(input)

		for item := range quantity {
			input <- item
		}
	}()

	for item := range discipline.Output() {
		_ = item
	}
}

// Here we model the worst case: when the quantity of operations for measuring time and
// calculating delays is equal to the quantity of operations for transmitting data
// items.
func BenchmarkDiscipline(b *testing.B) {
	quantity := b.N

	limit := Rate{
		Interval: time.Nanosecond,
		Quantity: 1,
	}

	input := make(chan int, quantity)

	opts := Opts[int]{
		Input: input,
		Limit: limit,
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	b.ResetTimer()

	go func() {
		defer close(input)

		for item := range quantity {
			input <- item
		}
	}()

	for item := range discipline.Output() {
		_ = item
	}
}
