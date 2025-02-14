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

func TestDisciplineTimeout(t *testing.T) {
	const (
		timeout = 100 * time.Millisecond
		pause   = 10 * timeout
	)

	data := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
	}

	exp1 := [][]int{
		{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10},
		{11}, {12}, {13}, {14}, {15}, {16}, {17}, {18}, {19}, {20},
		{21}, {22}, {23}, {24}, {25}, {26}, {27}, {28}, {29}, {30},
	}

	dur1At1 := []time.Duration{
		pause, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	dur1At2 := []time.Duration{
		0, pause, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	dur1At3 := []time.Duration{
		0, 0, pause, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	dur1At4 := []time.Duration{
		0, 0, 0, pause, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	dur1At5 := []time.Duration{
		0, 0, 0, 0, pause, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	dur1At6 := []time.Duration{
		0, 0, 0, 0, 0, pause, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	testDisciplineParallel(t, "1", data, 1, false, timeout, 1, pause, exp1, dur1At1)
	testDisciplineParallel(t, "1", data, 1, true, timeout, 1, pause, exp1, dur1At1)
	testDisciplineParallel(t, "1", data, 1, false, timeout, 2, pause, exp1, dur1At2)
	testDisciplineParallel(t, "1", data, 1, true, timeout, 2, pause, exp1, dur1At2)
	testDisciplineParallel(t, "1", data, 1, false, timeout, 3, pause, exp1, dur1At3)
	testDisciplineParallel(t, "1", data, 1, true, timeout, 3, pause, exp1, dur1At3)
	testDisciplineParallel(t, "1", data, 1, false, timeout, 4, pause, exp1, dur1At4)
	testDisciplineParallel(t, "1", data, 1, true, timeout, 4, pause, exp1, dur1At4)
	testDisciplineParallel(t, "1", data, 1, false, timeout, 5, pause, exp1, dur1At5)
	testDisciplineParallel(t, "1", data, 1, true, timeout, 5, pause, exp1, dur1At5)
	testDisciplineParallel(t, "1", data, 1, false, timeout, 6, pause, exp1, dur1At6)
	testDisciplineParallel(t, "1", data, 1, true, timeout, 6, pause, exp1, dur1At6)

	exp3At1 := [][]int{
		{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15},
		{16, 17, 18}, {19, 20, 21}, {22, 23, 24}, {25, 26, 27}, {28, 29, 30},
	}

	dur3At1 := []time.Duration{
		pause, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	exp3At2 := [][]int{
		{1}, {2, 3, 4}, {5, 6, 7}, {8, 9, 10}, {11, 12, 13}, {14, 15, 16},
		{17, 18, 19}, {20, 21, 22}, {23, 24, 25}, {26, 27, 28}, {29, 30},
	}

	dur3At2 := []time.Duration{
		timeout, pause - timeout, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	exp3At3 := [][]int{
		{1, 2}, {3, 4, 5}, {6, 7, 8}, {9, 10, 11}, {12, 13, 14}, {15, 16, 17},
		{18, 19, 20}, {21, 22, 23}, {24, 25, 26}, {27, 28, 29}, {30},
	}

	dur3At3 := []time.Duration{
		timeout, pause - timeout, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	exp3At4 := [][]int{
		{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15},
		{16, 17, 18}, {19, 20, 21}, {22, 23, 24}, {25, 26, 27}, {28, 29, 30},
	}

	dur3At4 := []time.Duration{
		0, pause, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	exp3At5 := [][]int{
		{1, 2, 3}, {4}, {5, 6, 7}, {8, 9, 10}, {11, 12, 13}, {14, 15, 16},
		{17, 18, 19}, {20, 21, 22}, {23, 24, 25}, {26, 27, 28}, {29, 30},
	}

	dur3At5 := []time.Duration{
		0, timeout, pause - timeout, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	exp3At6 := [][]int{
		{1, 2, 3}, {4, 5}, {6, 7, 8}, {9, 10, 11}, {12, 13, 14}, {15, 16, 17},
		{18, 19, 20}, {21, 22, 23}, {24, 25, 26}, {27, 28, 29}, {30},
	}

	dur3At6 := []time.Duration{
		0, timeout, pause - timeout, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	testDisciplineParallel(t, "3", data, 3, false, timeout, 1, pause, exp3At1, dur3At1)
	testDisciplineParallel(t, "3", data, 3, true, timeout, 1, pause, exp3At1, dur3At1)
	testDisciplineParallel(t, "3", data, 3, false, timeout, 2, pause, exp3At2, dur3At2)
	testDisciplineParallel(t, "3", data, 3, true, timeout, 2, pause, exp3At2, dur3At2)
	testDisciplineParallel(t, "3", data, 3, false, timeout, 3, pause, exp3At3, dur3At3)
	testDisciplineParallel(t, "3", data, 3, true, timeout, 3, pause, exp3At3, dur3At3)
	testDisciplineParallel(t, "3", data, 3, false, timeout, 4, pause, exp3At4, dur3At4)
	testDisciplineParallel(t, "3", data, 3, true, timeout, 4, pause, exp3At4, dur3At4)
	testDisciplineParallel(t, "3", data, 3, false, timeout, 5, pause, exp3At5, dur3At5)
	testDisciplineParallel(t, "3", data, 3, true, timeout, 5, pause, exp3At5, dur3At5)
	testDisciplineParallel(t, "3", data, 3, false, timeout, 6, pause, exp3At6, dur3At6)
	testDisciplineParallel(t, "3", data, 3, true, timeout, 6, pause, exp3At6, dur3At6)

	exp4At1 := [][]int{
		{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15, 16},
		{17, 18, 19, 20}, {21, 22, 23, 24}, {25, 26, 27, 28}, {29, 30},
	}

	dur4At1 := []time.Duration{
		pause, 0, 0, 0, 0, 0, 0, 0,
	}

	exp4At2 := [][]int{
		{1}, {2, 3, 4, 5}, {6, 7, 8, 9}, {10, 11, 12, 13}, {14, 15, 16, 17},
		{18, 19, 20, 21}, {22, 23, 24, 25}, {26, 27, 28, 29}, {30},
	}

	dur4At2 := []time.Duration{
		timeout, pause - timeout, 0, 0, 0, 0, 0, 0, 0,
	}

	exp4At3 := [][]int{
		{1, 2}, {3, 4, 5, 6}, {7, 8, 9, 10}, {11, 12, 13, 14}, {15, 16, 17, 18},
		{19, 20, 21, 22}, {23, 24, 25, 26}, {27, 28, 29, 30},
	}

	dur4At3 := []time.Duration{
		timeout, pause - timeout, 0, 0, 0, 0, 0, 0,
	}

	exp4At4 := [][]int{
		{1, 2, 3}, {4, 5, 6, 7}, {8, 9, 10, 11}, {12, 13, 14, 15}, {16, 17, 18, 19},
		{20, 21, 22, 23}, {24, 25, 26, 27}, {28, 29, 30},
	}

	dur4At4 := []time.Duration{
		timeout, pause - timeout, 0, 0, 0, 0, 0, 0,
	}

	exp4At5 := [][]int{
		{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15, 16},
		{17, 18, 19, 20}, {21, 22, 23, 24}, {25, 26, 27, 28}, {29, 30},
	}

	dur4At5 := []time.Duration{
		0, pause, 0, 0, 0, 0, 0, 0,
	}

	exp4At6 := [][]int{
		{1, 2, 3, 4}, {5}, {6, 7, 8, 9}, {10, 11, 12, 13}, {14, 15, 16, 17},
		{18, 19, 20, 21}, {22, 23, 24, 25}, {26, 27, 28, 29}, {30},
	}

	dur4At6 := []time.Duration{
		0, timeout, pause - timeout, 0, 0, 0, 0, 0, 0,
	}

	testDisciplineParallel(t, "4", data, 4, false, timeout, 1, pause, exp4At1, dur4At1)
	testDisciplineParallel(t, "4", data, 4, true, timeout, 1, pause, exp4At1, dur4At1)
	testDisciplineParallel(t, "4", data, 4, false, timeout, 2, pause, exp4At2, dur4At2)
	testDisciplineParallel(t, "4", data, 4, true, timeout, 2, pause, exp4At2, dur4At2)
	testDisciplineParallel(t, "4", data, 4, false, timeout, 3, pause, exp4At3, dur4At3)
	testDisciplineParallel(t, "4", data, 4, true, timeout, 3, pause, exp4At3, dur4At3)
	testDisciplineParallel(t, "4", data, 4, false, timeout, 4, pause, exp4At4, dur4At4)
	testDisciplineParallel(t, "4", data, 4, true, timeout, 4, pause, exp4At4, dur4At4)
	testDisciplineParallel(t, "4", data, 4, false, timeout, 5, pause, exp4At5, dur4At5)
	testDisciplineParallel(t, "4", data, 4, true, timeout, 5, pause, exp4At5, dur4At5)
	testDisciplineParallel(t, "4", data, 4, false, timeout, 6, pause, exp4At6, dur4At6)
	testDisciplineParallel(t, "4", data, 4, true, timeout, 6, pause, exp4At6, dur4At6)
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
	name string,
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
		name,
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
