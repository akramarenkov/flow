package join

import (
	"testing"
	"time"

	"github.com/akramarenkov/flow/join/internal/defs"
	"github.com/akramarenkov/flow/join/internal/inspect"
	"github.com/akramarenkov/safe"

	"github.com/stretchr/testify/require"
)

func TestOptsValidation(t *testing.T) {
	opts := Opts[int]{
		JoinSize: 10,
	}

	_, err := New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input: make(chan int),
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input:    make(chan int),
		JoinSize: 10,
		Timeout:  3 * time.Nanosecond,
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input:    make(chan int),
		JoinSize: 10,
		Timeout:  defs.MinTimeout,
	}

	_, err = New(opts)
	require.NoError(t, err)

	opts = Opts[int]{
		Input:    make(chan int),
		JoinSize: 10,
	}

	_, err = New(opts)
	require.NoError(t, err)
}

func TestDiscipline(t *testing.T) {
	for quantity := range safe.Iter[uint](100, 200) {
		for joinSize := range safe.Iter[uint](1, 20) {
			testDiscipline(t, quantity, joinSize, false, defs.TestTimeout)
			testDiscipline(t, quantity, joinSize, true, defs.TestTimeout)
			testDiscipline(t, quantity, joinSize, false, 0)
			testDiscipline(t, quantity, joinSize, true, 0)
		}
	}
}

func testDiscipline(
	t *testing.T,
	quantity uint,
	joinSize uint,
	noCopy bool,
	timeout time.Duration,
) {
	input := make(chan uint, joinSize)

	opts := Opts[uint]{
		Input:    input,
		JoinSize: joinSize,
		NoCopy:   noCopy,
		Timeout:  timeout,
	}

	discipline, err := New(opts)
	require.NoError(
		t,
		err,
		"quantity: %v, join size: %v, no copy: %v, timeout: %v",
		quantity,
		joinSize,
		noCopy,
		timeout,
	)

	inSequence := make([]uint, 0, quantity)
	outSequence := make([]uint, 0, quantity)

	expected := inspect.Expected(quantity, 1, joinSize)
	output := make([][]uint, 0, len(expected))

	go func() {
		defer close(input)

		for _, block := range inspect.Input(quantity, 1) {
			for _, item := range block {
				inSequence = append(inSequence, item)

				input <- item
			}
		}
	}()

	for join := range discipline.Output() {
		require.NotEmpty(
			t,
			join,
			"quantity: %v, join size: %v, no copy: %v, timeout: %v",
			quantity,
			joinSize,
			noCopy,
			timeout,
		)

		output = append(output, append([]uint(nil), join...))
		outSequence = append(outSequence, join...)

		if noCopy {
			discipline.Release()
		}
	}

	require.Equal(
		t,
		inSequence,
		outSequence,
		"quantity: %v, join size: %v, no copy: %v, timeout: %v",
		quantity,
		joinSize,
		noCopy,
		timeout,
	)

	require.Equal(
		t,
		expected,
		output,
		"quantity: %v, join size: %v, no copy: %v, timeout: %v",
		quantity,
		joinSize,
		noCopy,
		timeout,
	)
}

func TestDisciplineTimeout(t *testing.T) {
	for pauseAt := range safe.Iter[uint](50, 70) {
		t.Run(
			"",
			func(t *testing.T) {
				t.Parallel()
				testDisciplineTimeout(t, 100, 10, false, 500*time.Millisecond, pauseAt)
			},
		)

		t.Run(
			"",
			func(t *testing.T) {
				t.Parallel()
				testDisciplineTimeout(t, 100, 10, true, 500*time.Millisecond, pauseAt)
			},
		)
	}
}

func testDisciplineTimeout(
	t *testing.T,
	quantity uint,
	joinSize uint,
	noCopy bool,
	timeout time.Duration,
	pauseAt uint,
) {
	input := make(chan uint, joinSize)

	opts := Opts[uint]{
		Input:    input,
		JoinSize: joinSize,
		NoCopy:   noCopy,
		Timeout:  timeout,
	}

	require.NotZero(t, timeout)

	pauseAt = inspect.PickUpPauseAt(quantity, pauseAt, 1, joinSize)
	require.NotZero(
		t,
		pauseAt,
		"quantity: %v, join size: %v, no copy: %v, timeout: %v, pause at: %v",
		quantity,
		joinSize,
		noCopy,
		timeout,
		pauseAt,
	)

	pausetAtDuration := inspect.CalcPauseAtDuration(timeout)

	discipline, err := New(opts)
	require.NoError(
		t,
		err,
		"quantity: %v, join size: %v, no copy: %v, timeout: %v, pause at: %v",
		quantity,
		joinSize,
		noCopy,
		timeout,
		pauseAt,
	)

	inSequence := make([]uint, 0, quantity)
	outSequence := make([]uint, 0, quantity)

	expected := inspect.ExpectedWithTimeout(quantity, pauseAt, 1, joinSize)
	output := make([][]uint, 0, len(expected))

	go func() {
		defer close(input)

		for _, block := range inspect.Input(quantity, 1) {
			for _, item := range block {
				if item == pauseAt {
					time.Sleep(pausetAtDuration)
				}

				inSequence = append(inSequence, item)

				input <- item
			}
		}
	}()

	for join := range discipline.Output() {
		require.NotEmpty(
			t,
			join,
			"quantity: %v, join size: %v, no copy: %v, timeout: %v, pause at: %v",
			quantity,
			joinSize,
			noCopy,
			timeout,
			pauseAt,
		)

		output = append(output, append([]uint(nil), join...))
		outSequence = append(outSequence, join...)

		if noCopy {
			discipline.Release()
		}
	}

	require.Equal(
		t,
		inSequence,
		outSequence,
		"quantity: %v, join size: %v, no copy: %v, timeout: %v, pause at: %v",
		quantity,
		joinSize,
		noCopy,
		timeout,
		pauseAt,
	)

	require.Equal(
		t,
		expected,
		output,
		"quantity: %v, join size: %v, no copy: %v, timeout: %v, pause at: %v",
		quantity,
		joinSize,
		noCopy,
		timeout,
		pauseAt,
	)
}

func BenchmarkDiscipline(b *testing.B) {
	benchmarkDiscipline(b, 10, false, defs.TestTimeout, 1)
}

func BenchmarkDisciplineNoCopy(b *testing.B) {
	benchmarkDiscipline(b, 10, true, defs.TestTimeout, 1)
}

func BenchmarkDisciplineUntimeouted(b *testing.B) {
	benchmarkDiscipline(b, 10, false, 0, 1)
}

func BenchmarkDisciplineNoCopyUntimeouted(b *testing.B) {
	benchmarkDiscipline(b, 10, true, 0, 1)
}

func BenchmarkDisciplineInputCapacity0(b *testing.B) {
	benchmarkDiscipline(b, 10, false, 0, 0)
}

func BenchmarkDisciplineNoCopyInputCapacity0(b *testing.B) {
	benchmarkDiscipline(b, 10, true, 0, 0)
}

func BenchmarkDisciplineInputCapacity50(b *testing.B) {
	benchmarkDiscipline(b, 10, false, 0, 0.5)
}

func BenchmarkDisciplineNoCopyInputCapacity50(b *testing.B) {
	benchmarkDiscipline(b, 10, true, 0, 0.5)
}

func BenchmarkDisciplineInputCapacity100(b *testing.B) {
	benchmarkDiscipline(b, 10, false, 0, 1)
}

func BenchmarkDisciplineNoCopyInputCapacity100(b *testing.B) {
	benchmarkDiscipline(b, 10, true, 0, 1)
}

func BenchmarkDisciplineInputCapacity200(b *testing.B) {
	benchmarkDiscipline(b, 10, false, 0, 2)
}

func BenchmarkDisciplineNoCopyInputCapacity200(b *testing.B) {
	benchmarkDiscipline(b, 10, true, 0, 2)
}

func BenchmarkDisciplineInputCapacity300(b *testing.B) {
	benchmarkDiscipline(b, 10, false, 0, 3)
}

func BenchmarkDisciplineNoCopyInputCapacity300(b *testing.B) {
	benchmarkDiscipline(b, 10, true, 0, 3)
}

func BenchmarkDisciplineInputCapacity400(b *testing.B) {
	benchmarkDiscipline(b, 10, false, 0, 4)
}

func BenchmarkDisciplineNoCopyInputCapacity400(b *testing.B) {
	benchmarkDiscipline(b, 10, true, 0, 4)
}

func benchmarkDiscipline(
	b *testing.B,
	joinSize uint,
	noCopy bool,
	timeout time.Duration,
	inputCapFactor float64,
) {
	joinsQuantity, err := safe.IToI[uint](b.N)
	require.NoError(b, err)

	quantity := joinsQuantity * joinSize
	inputCap := uint(inputCapFactor * float64(joinSize))

	input := make(chan uint, inputCap)

	opts := Opts[uint]{
		Input:    input,
		JoinSize: joinSize,
		NoCopy:   noCopy,
		Timeout:  timeout,
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	b.ResetTimer()

	go func() {
		defer close(input)

		for item := uint(1); item <= quantity; item++ {
			input <- item
		}
	}()

	for range discipline.Output() {
		if noCopy {
			discipline.Release()
		}
	}
}