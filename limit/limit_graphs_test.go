package limit

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/akramarenkov/flow/internal/env"
	"github.com/akramarenkov/flow/limit/internal/research"

	"github.com/akramarenkov/safe"
	"github.com/akramarenkov/stressor"
	"github.com/go-echarts/go-echarts/v2/charts"
	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/require"
)

func TestGraphTicker(t *testing.T) {
	testGraphTicker(t, 1e2, time.Second, false)
	testGraphTicker(t, 1e2, 100*time.Millisecond, false)
	testGraphTicker(t, 1e3, 10*time.Millisecond, false)
	testGraphTicker(t, 1e3, 5*time.Millisecond, false)
	testGraphTicker(t, 1e3, time.Millisecond, false)
	testGraphTicker(t, 1e3, 100*time.Microsecond, false)

	testGraphTicker(t, 1e2, time.Second, true)
	testGraphTicker(t, 1e2, 100*time.Millisecond, true)
	testGraphTicker(t, 1e3, 10*time.Millisecond, true)
	testGraphTicker(t, 1e3, 5*time.Millisecond, true)
	testGraphTicker(t, 1e3, time.Millisecond, true)
	testGraphTicker(t, 1e3, 100*time.Microsecond, true)
}

func testGraphTicker(t *testing.T, quantity int, duration time.Duration, stress bool) {
	if os.Getenv(env.EnableGraphs) == "" {
		t.SkipNow()
	}

	if stress {
		stressor := stressor.New(stressor.Opts{})
		defer stressor.Stop()

		time.Sleep(time.Second)
	}

	times := make([]time.Duration, 0, quantity)

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	startedAt := time.Now()

	for range ticker.C {
		times = append(times, time.Since(startedAt))

		if len(times) == quantity {
			break
		}
	}

	createDelayerQuantitiesGraph(t, "Ticker", "ticker", times, duration, stress)
	createDelayerDeviationsGraph(t, "Ticker", "ticker", times, duration, stress)
}

func TestGraphSleep(t *testing.T) {
	testGraphSleep(t, 1e2, time.Second, false)
	testGraphSleep(t, 1e2, 100*time.Millisecond, false)
	testGraphSleep(t, 1e3, 10*time.Millisecond, false)
	testGraphSleep(t, 1e3, 5*time.Millisecond, false)
	testGraphSleep(t, 1e3, time.Millisecond, false)
	testGraphSleep(t, 1e3, 100*time.Microsecond, false)

	testGraphSleep(t, 1e2, time.Second, true)
	testGraphSleep(t, 1e2, 100*time.Millisecond, true)
	testGraphSleep(t, 1e3, 10*time.Millisecond, true)
	testGraphSleep(t, 1e3, 5*time.Millisecond, true)
	testGraphSleep(t, 1e3, time.Millisecond, true)
	testGraphSleep(t, 1e3, 100*time.Microsecond, true)
}

func testGraphSleep(t *testing.T, quantity int, duration time.Duration, stress bool) {
	if os.Getenv(env.EnableGraphs) == "" {
		t.SkipNow()
	}

	if stress {
		stressor := stressor.New(stressor.Opts{})
		defer stressor.Stop()

		time.Sleep(time.Second)
	}

	times := make([]time.Duration, quantity)

	startedAt := time.Now()

	for id := range quantity {
		time.Sleep(duration)

		times[id] = time.Since(startedAt)
	}

	createDelayerQuantitiesGraph(t, "Sleep", "sleep", times, duration, stress)
	createDelayerDeviationsGraph(t, "Sleep", "sleep", times, duration, stress)
}

func createDelayerQuantitiesGraph(
	t *testing.T,
	titlePerfix string,
	fileNamePerfix string,
	times []time.Duration,
	delay time.Duration,
	stress bool,
) {
	quantities, interval := research.QuantityPerInterval(times, 0, delay)

	axisY, axisX := research.QotToBarChart(quantities)

	expectedDuration := time.Duration(len(times)) * delay

	subtitleAdd := fmt.Sprintf(
		"duration: %s, %s",
		delay,
		fmtTotalDuration(expectedDuration, times),
	)

	fileNameAdd := fmt.Sprintf("%s_quantities_duration_%s", fileNamePerfix, delay)

	createGraph(
		t,
		titlePerfix+" quantities over time",
		subtitleAdd,
		fileNameAdd,
		"quantities",
		len(times),
		interval.String(),
		stress,
		axisY,
		axisX,
	)
}

func createDelayerDeviationsGraph(
	t *testing.T,
	titlePerfix string,
	fileNamePerfix string,
	times []time.Duration,
	duration time.Duration,
	stress bool,
) {
	deviations := research.Deviations(times, duration)

	axisY, axisX := research.DeviationsToBarChart(deviations)

	subtitleAdd := fmt.Sprintf("duration: %s", duration)
	fileNameAdd := fmt.Sprintf("%s_deviations_duration_%s", fileNamePerfix, duration)

	createGraph(
		t,
		titlePerfix+" deviations",
		subtitleAdd,
		fileNameAdd,
		"deviations",
		len(times),
		"1%",
		stress,
		axisY,
		axisX,
	)
}

func TestGraphDiscipline(t *testing.T) {
	testGraphDiscipline(
		t,
		1e4+1,
		Rate{
			Interval: time.Second,
			Quantity: 1e3,
		},
		false,
	)

	testGraphDiscipline(
		t,
		1e4+1,
		Rate{
			Interval: 100 * time.Millisecond,
			Quantity: 1e2,
		},
		false,
	)

	testGraphDiscipline(
		t,
		1e4+1,
		Rate{
			Interval: 10 * time.Millisecond,
			Quantity: 1e1,
		},
		false,
	)

	testGraphDiscipline(
		t,
		1e5+1,
		Rate{
			Interval: time.Millisecond,
			Quantity: 1e1,
		},
		false,
	)

	testGraphDiscipline(
		t,
		1e4+1,
		Rate{
			Interval: time.Nanosecond,
			Quantity: 1,
		},
		false,
	)

	testGraphDiscipline(
		t,
		1e4+1,
		Rate{
			Interval: time.Second,
			Quantity: 1e3,
		},
		true,
	)

	testGraphDiscipline(
		t,
		1e4+1,
		Rate{
			Interval: 100 * time.Millisecond,
			Quantity: 1e2,
		},
		true,
	)

	testGraphDiscipline(
		t,
		1e4+1,
		Rate{
			Interval: 10 * time.Millisecond,
			Quantity: 1e1,
		},
		true,
	)

	testGraphDiscipline(
		t,
		1e5+1,
		Rate{
			Interval: time.Millisecond,
			Quantity: 1e1,
		},
		true,
	)

	testGraphDiscipline(
		t,
		1e4+1,
		Rate{
			Interval: time.Nanosecond,
			Quantity: 1,
		},
		true,
	)
}

func testGraphDiscipline(t *testing.T, quantity int, limit Rate, stress bool) {
	if os.Getenv(env.EnableGraphs) == "" {
		t.SkipNow()
	}

	if stress {
		stressor := stressor.New(stressor.Opts{})
		defer stressor.Stop()

		time.Sleep(time.Second)
	}

	input := make(chan int, quantity)

	opts := Opts[int]{
		Input: input,
		Limit: limit,
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	times := make([]time.Duration, 0, quantity)

	startedAt := time.Now()

	go func() {
		defer close(input)

		for item := range quantity {
			input <- item
		}
	}()

	for range discipline.Output() {
		times = append(times, time.Since(startedAt))
	}

	createQuantitiesGraph(t, times, limit, stress)
}

func createQuantitiesGraph(
	t *testing.T,
	times []time.Duration,
	limit Rate,
	stress bool,
) {
	quantities, interval := research.QuantityPerInterval(times, 100, 0)

	axisY, axisX := research.QotToBarChart(quantities)

	limitQuantity, err := safe.IToI[time.Duration](limit.Quantity)
	require.NoError(t, err)

	expectedDuration := time.Duration(len(times)) * limit.Interval / limitQuantity

	subtitleAdd := fmt.Sprintf(
		"limit: {quantity: %d, interval: %s}, %s",
		limit.Quantity,
		limit.Interval,
		fmtTotalDuration(expectedDuration, times),
	)

	fileNameAdd := fmt.Sprintf(
		"quantities_limit_quantity_%d_limit_interval_%s",
		limit.Quantity,
		limit.Interval,
	)

	createGraph(
		t,
		"Quantities over time",
		subtitleAdd,
		fileNameAdd,
		"quantities",
		len(times),
		interval.String(),
		stress,
		axisY,
		axisX,
	)
}

func createGraph(
	t *testing.T,
	title string,
	subtitleAdd string,
	fileNameAdd string,
	seriesName string,
	totalQuantity int,
	graphInterval string,
	stress bool,
	series []chartsopts.BarData,
	abscissa interface{},
) {
	subtitle := fmt.Sprintf(
		"Total quantity: %d, graph interval: %s, %s, stress system: %t, time: %s",
		totalQuantity,
		graphInterval,
		subtitleAdd,
		stress,
		time.Now().Format(time.RFC3339),
	)

	fileName := fmt.Sprintf(
		"graph_%d_%s_stress_%t.html",
		totalQuantity,
		fileNameAdd,
		stress,
	)

	createBarGraph(
		t,
		title,
		subtitle,
		fileName,
		seriesName,
		series,
		abscissa,
	)
}

func createBarGraph(
	t *testing.T,
	title string,
	subtitle string,
	fileName string,
	seriesName string,
	series []chartsopts.BarData,
	abscissa interface{},
) {
	if len(series) == 0 {
		return
	}

	chart := charts.NewBar()

	chart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    title,
				Subtitle: subtitle,
			},
		),
	)

	chart.SetXAxis(abscissa).AddSeries(seriesName, series)

	file, err := os.Create(fileName)
	require.NoError(t, err)

	defer file.Close()

	err = chart.Render(file)
	require.NoError(t, err)
}

func fmtTotalDuration(expected time.Duration, times []time.Duration) string {
	formatted := fmt.Sprintf(
		"total duration: {expected:  %s, actual: %s}",
		expected,
		research.TotalDuration(times),
	)

	return formatted
}
