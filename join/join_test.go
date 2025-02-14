package join

import (
	"testing"
	"time"

	"github.com/akramarenkov/flow/join/internal/defaults"

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
	}

	_, err = New(opts)
	require.NoError(t, err)

	opts = Opts[int]{
		Input:    make(chan int),
		JoinSize: 10,
		Timeout:  defaults.TestTimeout,
	}

	_, err = New(opts)
	require.NoError(t, err)

	opts = Opts[int]{
		Input:    make(chan int),
		JoinSize: 10,
		Timeout:  -defaults.TestTimeout,
	}

	_, err = New(opts)
	require.NoError(t, err)
}

func TestDiscipline(t *testing.T) {
	data := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
	}

	expected1 := [][]int{
		{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10},
		{11}, {12}, {13}, {14}, {15}, {16}, {17}, {18}, {19}, {20},
		{21}, {22}, {23}, {24}, {25}, {26}, {27}, {28}, {29}, {30},
	}

	expected3 := [][]int{
		{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15},
		{16, 17, 18}, {19, 20, 21}, {22, 23, 24}, {25, 26, 27}, {28, 29, 30},
	}

	expected4 := [][]int{
		{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15, 16},
		{17, 18, 19, 20}, {21, 22, 23, 24}, {25, 26, 27, 28}, {29, 30},
	}

	testDiscipline(t, data, 1, false, defaults.TestTimeout, 0, 0, expected1, nil)
	testDiscipline(t, data, 1, true, defaults.TestTimeout, 0, 0, expected1, nil)
	testDiscipline(t, data, 1, false, 0, 0, 0, expected1, nil)
	testDiscipline(t, data, 1, true, 0, 0, 0, expected1, nil)

	testDiscipline(t, data, 3, false, defaults.TestTimeout, 0, 0, expected3, nil)
	testDiscipline(t, data, 3, true, defaults.TestTimeout, 0, 0, expected3, nil)
	testDiscipline(t, data, 3, false, 0, 0, 0, expected3, nil)
	testDiscipline(t, data, 3, true, 0, 0, 0, expected3, nil)

	testDiscipline(t, data, 4, false, defaults.TestTimeout, 0, 0, expected4, nil)
	testDiscipline(t, data, 4, true, defaults.TestTimeout, 0, 0, expected4, nil)
	testDiscipline(t, data, 4, false, 0, 0, 0, expected4, nil)
	testDiscipline(t, data, 4, true, 0, 0, 0, expected4, nil)
}

func TestDisciplineTimeoutJoinSize1(t *testing.T) {
	const (
		timeout = 100 * time.Millisecond
		pause   = 10 * timeout
	)

	data := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
	}

	exp := [][]int{
		{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10},
		{11}, {12}, {13}, {14}, {15}, {16}, {17}, {18}, {19}, {20},
		{21}, {22}, {23}, {24}, {25}, {26}, {27}, {28}, {29}, {30},
	}

	durAt1 := []time.Duration{
		pause, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	durAt2 := []time.Duration{
		0, pause, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	durAt3 := []time.Duration{
		0, 0, pause, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	durAt4 := []time.Duration{
		0, 0, 0, pause, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	durAt5 := []time.Duration{
		0, 0, 0, 0, pause, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	durAt6 := []time.Duration{
		0, 0, 0, 0, 0, pause, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	testDisciplineParallel(t, data, 1, false, timeout, 1, pause, exp, durAt1)
	testDisciplineParallel(t, data, 1, true, timeout, 1, pause, exp, durAt1)
	testDisciplineParallel(t, data, 1, false, timeout, 2, pause, exp, durAt2)
	testDisciplineParallel(t, data, 1, true, timeout, 2, pause, exp, durAt2)
	testDisciplineParallel(t, data, 1, false, timeout, 3, pause, exp, durAt3)
	testDisciplineParallel(t, data, 1, true, timeout, 3, pause, exp, durAt3)
	testDisciplineParallel(t, data, 1, false, timeout, 4, pause, exp, durAt4)
	testDisciplineParallel(t, data, 1, true, timeout, 4, pause, exp, durAt4)
	testDisciplineParallel(t, data, 1, false, timeout, 5, pause, exp, durAt5)
	testDisciplineParallel(t, data, 1, true, timeout, 5, pause, exp, durAt5)
	testDisciplineParallel(t, data, 1, false, timeout, 6, pause, exp, durAt6)
	testDisciplineParallel(t, data, 1, true, timeout, 6, pause, exp, durAt6)
}

func TestDisciplineTimeoutJoinSize3(t *testing.T) {
	const (
		timeout = 200 * time.Millisecond
		pause   = 10 * timeout
	)

	data := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
	}

	expAt1 := [][]int{
		{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15},
		{16, 17, 18}, {19, 20, 21}, {22, 23, 24}, {25, 26, 27}, {28, 29, 30},
	}

	durAt1 := []time.Duration{
		pause, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	expAt2 := [][]int{
		{1}, {2, 3, 4}, {5, 6, 7}, {8, 9, 10}, {11, 12, 13}, {14, 15, 16},
		{17, 18, 19}, {20, 21, 22}, {23, 24, 25}, {26, 27, 28}, {29, 30},
	}

	durAt2 := []time.Duration{
		timeout, pause - timeout, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	expAt3 := [][]int{
		{1, 2}, {3, 4, 5}, {6, 7, 8}, {9, 10, 11}, {12, 13, 14}, {15, 16, 17},
		{18, 19, 20}, {21, 22, 23}, {24, 25, 26}, {27, 28, 29}, {30},
	}

	durAt3 := []time.Duration{
		timeout, pause - timeout, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	expAt4 := [][]int{
		{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15},
		{16, 17, 18}, {19, 20, 21}, {22, 23, 24}, {25, 26, 27}, {28, 29, 30},
	}

	durAt4 := []time.Duration{
		0, pause, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	expAt5 := [][]int{
		{1, 2, 3}, {4}, {5, 6, 7}, {8, 9, 10}, {11, 12, 13}, {14, 15, 16},
		{17, 18, 19}, {20, 21, 22}, {23, 24, 25}, {26, 27, 28}, {29, 30},
	}

	durAt5 := []time.Duration{
		0, timeout, pause - timeout, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	expAt6 := [][]int{
		{1, 2, 3}, {4, 5}, {6, 7, 8}, {9, 10, 11}, {12, 13, 14}, {15, 16, 17},
		{18, 19, 20}, {21, 22, 23}, {24, 25, 26}, {27, 28, 29}, {30},
	}

	durAt6 := []time.Duration{
		0, timeout, pause - timeout, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	testDisciplineParallel(t, data, 3, false, timeout, 1, pause, expAt1, durAt1)
	testDisciplineParallel(t, data, 3, true, timeout, 1, pause, expAt1, durAt1)
	testDisciplineParallel(t, data, 3, false, timeout, 2, pause, expAt2, durAt2)
	testDisciplineParallel(t, data, 3, true, timeout, 2, pause, expAt2, durAt2)
	testDisciplineParallel(t, data, 3, false, timeout, 3, pause, expAt3, durAt3)
	testDisciplineParallel(t, data, 3, true, timeout, 3, pause, expAt3, durAt3)
	testDisciplineParallel(t, data, 3, false, timeout, 4, pause, expAt4, durAt4)
	testDisciplineParallel(t, data, 3, true, timeout, 4, pause, expAt4, durAt4)
	testDisciplineParallel(t, data, 3, false, timeout, 5, pause, expAt5, durAt5)
	testDisciplineParallel(t, data, 3, true, timeout, 5, pause, expAt5, durAt5)
	testDisciplineParallel(t, data, 3, false, timeout, 6, pause, expAt6, durAt6)
	testDisciplineParallel(t, data, 3, true, timeout, 6, pause, expAt6, durAt6)
}

func TestDisciplineTimeoutJoinSize4(t *testing.T) {
	const (
		timeout = 300 * time.Millisecond
		pause   = 10 * timeout
	)

	data := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
	}

	expAt1 := [][]int{
		{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15, 16},
		{17, 18, 19, 20}, {21, 22, 23, 24}, {25, 26, 27, 28}, {29, 30},
	}

	durAt1 := []time.Duration{
		pause, 0, 0, 0, 0, 0, 0, 0,
	}

	expAt2 := [][]int{
		{1}, {2, 3, 4, 5}, {6, 7, 8, 9}, {10, 11, 12, 13}, {14, 15, 16, 17},
		{18, 19, 20, 21}, {22, 23, 24, 25}, {26, 27, 28, 29}, {30},
	}

	durAt2 := []time.Duration{
		timeout, pause - timeout, 0, 0, 0, 0, 0, 0, 0,
	}

	expAt3 := [][]int{
		{1, 2}, {3, 4, 5, 6}, {7, 8, 9, 10}, {11, 12, 13, 14}, {15, 16, 17, 18},
		{19, 20, 21, 22}, {23, 24, 25, 26}, {27, 28, 29, 30},
	}

	durAt3 := []time.Duration{
		timeout, pause - timeout, 0, 0, 0, 0, 0, 0,
	}

	expAt4 := [][]int{
		{1, 2, 3}, {4, 5, 6, 7}, {8, 9, 10, 11}, {12, 13, 14, 15}, {16, 17, 18, 19},
		{20, 21, 22, 23}, {24, 25, 26, 27}, {28, 29, 30},
	}

	durAt4 := []time.Duration{
		timeout, pause - timeout, 0, 0, 0, 0, 0, 0,
	}

	expAt5 := [][]int{
		{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15, 16},
		{17, 18, 19, 20}, {21, 22, 23, 24}, {25, 26, 27, 28}, {29, 30},
	}

	durAt5 := []time.Duration{
		0, pause, 0, 0, 0, 0, 0, 0,
	}

	expAt6 := [][]int{
		{1, 2, 3, 4}, {5}, {6, 7, 8, 9}, {10, 11, 12, 13}, {14, 15, 16, 17},
		{18, 19, 20, 21}, {22, 23, 24, 25}, {26, 27, 28, 29}, {30},
	}

	durAt6 := []time.Duration{
		0, timeout, pause - timeout, 0, 0, 0, 0, 0, 0,
	}

	testDisciplineParallel(t, data, 4, false, timeout, 1, pause, expAt1, durAt1)
	testDisciplineParallel(t, data, 4, true, timeout, 1, pause, expAt1, durAt1)
	testDisciplineParallel(t, data, 4, false, timeout, 2, pause, expAt2, durAt2)
	testDisciplineParallel(t, data, 4, true, timeout, 2, pause, expAt2, durAt2)
	testDisciplineParallel(t, data, 4, false, timeout, 3, pause, expAt3, durAt3)
	testDisciplineParallel(t, data, 4, true, timeout, 3, pause, expAt3, durAt3)
	testDisciplineParallel(t, data, 4, false, timeout, 4, pause, expAt4, durAt4)
	testDisciplineParallel(t, data, 4, true, timeout, 4, pause, expAt4, durAt4)
	testDisciplineParallel(t, data, 4, false, timeout, 5, pause, expAt5, durAt5)
	testDisciplineParallel(t, data, 4, true, timeout, 5, pause, expAt5, durAt5)
	testDisciplineParallel(t, data, 4, false, timeout, 6, pause, expAt6, durAt6)
	testDisciplineParallel(t, data, 4, true, timeout, 6, pause, expAt6, durAt6)
}

func testDiscipline(
	t *testing.T,
	data []int,
	joinSize uint,
	noCopy bool,
	timeout time.Duration,
	pauseAt int,
	pauseDuration time.Duration,
	expected [][]int,
	expectedDurations []time.Duration,
) {
	input := make(chan int, joinSize)

	opts := Opts[int]{
		Input:    input,
		JoinSize: joinSize,
		NoCopy:   noCopy,
		Timeout:  timeout,
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	durations := make([]time.Duration, 0, len(expected))
	output := make([][]int, 0, len(expected))

	go func() {
		defer close(input)

		for _, item := range data {
			if item == pauseAt {
				time.Sleep(pauseDuration)
			}

			input <- item
		}
	}()

	previous := time.Now()

	for join := range discipline.Output() {
		durations = append(durations, time.Since(previous))

		if noCopy {
			output = append(output, append([]int(nil), join...))
		} else {
			output = append(output, join)
		}

		discipline.Release()

		previous = time.Now()
	}

	require.Equal(t, expected, output)

	if len(expectedDurations) != 0 {
		require.Len(t, durations, len(expectedDurations))
	}

	for id, expected := range expectedDurations {
		if expected == 0 {
			require.Less(t, durations[id], timeout)
			continue
		}

		require.InEpsilon(t, expected, durations[id], 0.05)
	}
}

func testDisciplineParallel(
	t *testing.T,
	data []int,
	joinSize uint,
	noCopy bool,
	timeout time.Duration,
	pauseAt int,
	pauseDuration time.Duration,
	expected [][]int,
	expectedDurations []time.Duration,
) {
	t.Run(
		"",
		func(t *testing.T) {
			t.Parallel()

			testDiscipline(
				t,
				data,
				joinSize,
				noCopy,
				timeout,
				pauseAt,
				pauseDuration,
				expected,
				expectedDurations,
			)
		},
	)
}

func BenchmarkDiscipline(b *testing.B) {
	benchmarkDiscipline(b, false, defaults.TestTimeout, 1)
}

func BenchmarkDisciplineNoCopy(b *testing.B) {
	benchmarkDiscipline(b, true, defaults.TestTimeout, 1)
}

func BenchmarkDisciplineUntimeouted(b *testing.B) {
	benchmarkDiscipline(b, false, 0, 1)
}

func BenchmarkDisciplineNoCopyUntimeouted(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 1)
}

func BenchmarkDisciplineInputCapacity0(b *testing.B) {
	benchmarkDiscipline(b, false, 0, 0)
}

func BenchmarkDisciplineNoCopyInputCapacity0(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0)
}

func BenchmarkDisciplineInputCapacity50(b *testing.B) {
	benchmarkDiscipline(b, false, 0, 0.5)
}

func BenchmarkDisciplineNoCopyInputCapacity50(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0.5)
}

func BenchmarkDisciplineInputCapacity100(b *testing.B) {
	benchmarkDiscipline(b, false, 0, 1)
}

func BenchmarkDisciplineNoCopyInputCapacity100(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 1)
}

func BenchmarkDisciplineInputCapacity200(b *testing.B) {
	benchmarkDiscipline(b, false, 0, 2)
}

func BenchmarkDisciplineNoCopyInputCapacity200(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 2)
}

func BenchmarkDisciplineInputCapacity300(b *testing.B) {
	benchmarkDiscipline(b, false, 0, 3)
}

func BenchmarkDisciplineNoCopyInputCapacity300(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 3)
}

func BenchmarkDisciplineInputCapacity400(b *testing.B) {
	benchmarkDiscipline(b, false, 0, 4)
}

func BenchmarkDisciplineNoCopyInputCapacity400(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 4)
}

func benchmarkDiscipline(
	b *testing.B,
	noCopy bool,
	timeout time.Duration,
	capacityFactor float64,
) {
	const joinSize = 10

	joinsQuantity := b.N

	sizeOfJoin, err := safe.IToI[int](joinSize)
	require.NoError(b, err)

	quantity := joinsQuantity * sizeOfJoin

	input := make(chan int, int(capacityFactor*float64(joinSize)))

	opts := Opts[int]{
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

		for item := range quantity {
			input <- item
		}
	}()

	for range discipline.Output() {
		discipline.Release()
	}
}
