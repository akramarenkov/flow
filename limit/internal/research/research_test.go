package research

import (
	"testing"
	"time"

	"github.com/akramarenkov/flow/internal/qot"
	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/require"
)

func TestQuantityPerIntervalSplitByInterval(t *testing.T) {
	times := []time.Duration{
		0,
		time.Millisecond,
		2 * time.Millisecond,
		5 * time.Millisecond,
		9 * time.Millisecond,
		11 * time.Millisecond,
		13 * time.Millisecond,
		17 * time.Millisecond,
	}

	interval := 10 * time.Millisecond

	expected := []qot.QoT{
		{
			Quantity: 5,
			Time:     0,
		},
		{
			Quantity: 3,
			Time:     10 * time.Millisecond,
		},
	}

	expectedAxisY := []chartsopts.BarData{
		{
			Name:  "0s",
			Value: uint(5),
			Tooltip: &chartsopts.Tooltip{
				Show: chartsopts.Bool(true),
			},
		},
		{
			Name:  "10ms",
			Value: uint(3),
			Tooltip: &chartsopts.Tooltip{
				Show: chartsopts.Bool(true),
			},
		},
	}

	expectedAxisX := []int{
		0,
		1,
	}

	quantities, actualInterval := QuantityPerInterval(times, 0, interval)
	require.Equal(t, expected, quantities)
	require.Equal(t, interval, actualInterval)

	axisY, axisX := QotToBarChart(quantities)
	require.Equal(t, expectedAxisY, axisY)
	require.Equal(t, expectedAxisX, axisX)
}

func TestQuantityPerIntervalSplitByIntervalEntirely(t *testing.T) {
	times := []time.Duration{
		0,
		time.Millisecond,
		2 * time.Millisecond,
		5 * time.Millisecond,
		9 * time.Millisecond,
		11 * time.Millisecond,
		13 * time.Millisecond,
		20 * time.Millisecond,
	}

	interval := 10 * time.Millisecond

	expected := []qot.QoT{
		{
			Quantity: 5,
			Time:     0,
		},
		{
			Quantity: 2,
			Time:     10 * time.Millisecond,
		},
		{
			Quantity: 1,
			Time:     20 * time.Millisecond,
		},
	}

	quantities, actualInterval := QuantityPerInterval(times, 0, interval)
	require.Equal(t, expected, quantities)
	require.Equal(t, interval, actualInterval)
}

func TestQuantityPerIntervalSplitByIntervalsNumber(t *testing.T) {
	times := []time.Duration{
		0,
		time.Millisecond,
		2 * time.Millisecond,
		5 * time.Millisecond,
		9 * time.Millisecond,
		11 * time.Millisecond,
		13 * time.Millisecond,
		17 * time.Millisecond,
	}

	intervalsNumber := 2

	expectedInterval := 8*time.Millisecond + 500*time.Microsecond + time.Nanosecond

	expected := []qot.QoT{
		{
			Quantity: 4,
			Time:     0,
		},
		{
			Quantity: 4,
			Time:     expectedInterval,
		},
	}

	quantities, actualInterval := QuantityPerInterval(times, intervalsNumber, 0)
	require.Equal(t, expected, quantities)
	require.Equal(t, expectedInterval, actualInterval)
}

func TestQuantityPerIntervalZeroInput(t *testing.T) {
	quantities, actualInterval := QuantityPerInterval(nil, 0, time.Second)
	require.Equal(t, []qot.QoT(nil), quantities)
	require.Equal(t, time.Duration(0), actualInterval)

	quantities, actualInterval = QuantityPerInterval([]time.Duration{}, 0, time.Second)
	require.Equal(t, []qot.QoT(nil), quantities)
	require.Equal(t, time.Duration(0), actualInterval)
}

func TestQuantityPerIntervalZeroSplit(t *testing.T) {
	quantities, actualInterval := QuantityPerInterval([]time.Duration{1, 2}, 0, 0)
	require.Equal(t, []qot.QoT(nil), quantities)
	require.Equal(t, time.Duration(0), actualInterval)
}

func TestQuantityPerIntervalSmallRatio(t *testing.T) {
	times := []time.Duration{
		0,
		time.Nanosecond,
		2 * time.Nanosecond,
		5 * time.Nanosecond,
	}

	intervalsNumber := 10

	expectedInterval := time.Nanosecond

	expected := []qot.QoT{
		{
			Quantity: 1,
			Time:     0,
		},
		{
			Quantity: 1,
			Time:     1,
		},
		{
			Quantity: 1,
			Time:     2,
		},
		{
			Quantity: 0,
			Time:     3,
		},
		{
			Quantity: 0,
			Time:     4,
		},
		{
			Quantity: 1,
			Time:     5,
		},
		{
			Quantity: 0,
			Time:     6,
		},
		{
			Quantity: 0,
			Time:     7,
		},
		{
			Quantity: 0,
			Time:     8,
		},
		{
			Quantity: 0,
			Time:     9,
		},
	}

	quantities, actualInterval := QuantityPerInterval(times, intervalsNumber, 0)
	require.Equal(t, expected, quantities)
	require.Equal(t, expectedInterval, actualInterval)
}

func TestDeviations(t *testing.T) {
	times := []time.Duration{
		-100 * time.Microsecond,
		900 * time.Microsecond,
		900 * time.Microsecond,
		2000 * time.Microsecond,
		2800 * time.Microsecond,
		3700 * time.Microsecond,
		4700 * time.Microsecond,
		5500 * time.Microsecond,
		6600 * time.Microsecond,
		7600 * time.Microsecond,
		8600 * time.Microsecond,
		9600 * time.Microsecond,
		10700 * time.Microsecond,
		12700 * time.Microsecond,
		14800 * time.Microsecond,
	}

	expected := make(map[int]int, 201)

	for deviation := -100; deviation <= 100; deviation++ {
		expected[deviation] = 0
	}

	expected[-100] = 2
	expected[-20] = 2
	expected[-10] = 1
	expected[0] = 5
	expected[10] = 3
	expected[100] = 2

	deviations := Deviations(times, time.Millisecond)
	require.Equal(t, expected, deviations)
}

func TestDeviationsZeroInput(t *testing.T) {
	deviations := Deviations(nil, time.Millisecond)
	require.Equal(t, map[int]int(nil), deviations)

	deviations = Deviations([]time.Duration{}, time.Millisecond)
	require.Equal(t, map[int]int(nil), deviations)
}

func TestConvertRelativeDeviationsToBarEcharts(t *testing.T) {
	deviations := map[int]int{
		-10: 2,
		0:   10,
		10:  3,
		80:  1,
	}

	expectedAxisY := []chartsopts.BarData{
		{
			Name:  "-10%",
			Value: 2,
			Tooltip: &chartsopts.Tooltip{
				Show: chartsopts.Bool(true),
			},
		},
		{
			Name:  "0%",
			Value: 10,
			Tooltip: &chartsopts.Tooltip{
				Show: chartsopts.Bool(true),
			},
		},
		{
			Name:  "10%",
			Value: 3,
			Tooltip: &chartsopts.Tooltip{
				Show: chartsopts.Bool(true),
			},
		},
		{
			Name:  "80%",
			Value: 1,
			Tooltip: &chartsopts.Tooltip{
				Show: chartsopts.Bool(true),
			},
		},
	}

	expectedAxisX := []int{
		-10,
		0,
		10,
		80,
	}

	axisY, axisX := DeviationsToBarChart(deviations)
	require.Equal(t, expectedAxisY, axisY)
	require.Equal(t, expectedAxisX, axisX)
}

func TestTotalDuration(t *testing.T) {
	durations := []time.Duration{
		17 * time.Millisecond,
		time.Millisecond,
		2 * time.Millisecond,
		9 * time.Millisecond,
		5 * time.Millisecond,
		11 * time.Millisecond,
		13 * time.Millisecond,
		0,
	}

	require.Equal(t, time.Duration(0), TotalDuration(nil))
	require.Equal(t, time.Duration(0), TotalDuration([]time.Duration{}))
	require.Equal(t, 17*time.Millisecond, TotalDuration(durations))
}
