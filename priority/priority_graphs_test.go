package priority

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/akramarenkov/flow/internal/consts"
	"github.com/akramarenkov/flow/internal/env"
	"github.com/akramarenkov/flow/priority/divider"
	"github.com/akramarenkov/flow/priority/internal/measuring"
	"github.com/akramarenkov/flow/priority/internal/research"
	"github.com/akramarenkov/flow/priority/internal/unmanaged"

	"github.com/go-echarts/go-echarts/v2/charts"
	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/require"
)

func TestGraphFairEven(t *testing.T) {
	testGraphFairEven(t, 1, false)
	testGraphFairEven(t, 10, false)
	testGraphFairEven(t, 100, false)
}

func TestGraphFairEvenUnbuffered(t *testing.T) {
	testGraphFairEven(t, 1, true)
	testGraphFairEven(t, 10, true)
	testGraphFairEven(t, 100, true)
}

func testGraphFairEven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv(env.EnableGraphs) == "" {
		t.SkipNow()
	}

	msr, err := measuring.NewMeasurer(6*factor, inputCapacity(unbufferedInput)...)
	require.NoError(t, err)

	msr.AddWrite(1, 4000*factor)

	msr.AddWrite(2, 500*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 500*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 4*time.Second)
	msr.AddWrite(2, 1000*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 2000*factor)

	msr.AddWrite(3, 500*factor)
	msr.AddWaitDevastation(3)
	msr.AddDelay(3, 5*time.Second)
	msr.AddWrite(3, 4000*factor)

	msr.SetProcessingDuration(1, 10*time.Millisecond)
	msr.SetProcessingDuration(2, 10*time.Millisecond)
	msr.SetProcessingDuration(3, 10*time.Millisecond)

	opts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measurements, err := msr.Play(discipline)
	require.NoError(t, err)

	createGraphs(
		t,
		"Fair divider, even time processing",
		"fair_even",
		opts.HandlersQuantity,
		unbufferedInput,
		measurements,
	)
}

func TestGraphFairUneven(t *testing.T) {
	testGraphFairUneven(t, 1, false)
	testGraphFairUneven(t, 10, false)
	testGraphFairUneven(t, 100, false)
}

func TestGraphFairUnevenUnbuffered(t *testing.T) {
	testGraphFairUneven(t, 1, true)
	testGraphFairUneven(t, 10, true)
	testGraphFairUneven(t, 100, true)
}

func testGraphFairUneven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv(env.EnableGraphs) == "" {
		t.SkipNow()
	}

	msr, err := measuring.NewMeasurer(6*factor, inputCapacity(unbufferedInput)...)
	require.NoError(t, err)

	msr.AddWrite(1, 450*factor)

	msr.AddWrite(2, 100*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 100*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 4*time.Second)
	msr.AddWrite(2, 200*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 400*factor)

	msr.AddWrite(3, 500*factor)
	msr.AddWaitDevastation(3)
	msr.AddDelay(3, 6*time.Second)
	msr.AddWrite(3, 3000*factor)

	msr.SetProcessingDuration(1, 100*time.Millisecond)
	msr.SetProcessingDuration(2, 50*time.Millisecond)
	msr.SetProcessingDuration(3, 10*time.Millisecond)

	opts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measurements, err := msr.Play(discipline)
	require.NoError(t, err)

	createGraphs(
		t,
		"Fair divider, uneven time processing",
		"fair_uneven",
		opts.HandlersQuantity,
		unbufferedInput,
		measurements,
	)
}

func TestGraphRateEven(t *testing.T) {
	testGraphRateEven(t, 1, false)
	testGraphRateEven(t, 10, false)
	testGraphRateEven(t, 100, false)
}

func TestGraphRateEvenUnbuffered(t *testing.T) {
	testGraphRateEven(t, 1, true)
	testGraphRateEven(t, 10, true)
	testGraphRateEven(t, 100, true)
}

func testGraphRateEven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv(env.EnableGraphs) == "" {
		t.SkipNow()
	}

	msr, err := measuring.NewMeasurer(6*factor, inputCapacity(unbufferedInput)...)
	require.NoError(t, err)

	msr.AddWrite(1, 3800*factor)

	msr.AddWrite(2, 1500*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 750*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 4*time.Second)
	msr.AddWrite(2, 700*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 1800*factor)

	msr.AddWrite(3, 1000*factor)
	msr.AddWaitDevastation(3)
	msr.AddDelay(3, 8*time.Second)
	msr.AddWrite(3, 4000*factor)

	msr.SetProcessingDuration(1, 10*time.Millisecond)
	msr.SetProcessingDuration(2, 10*time.Millisecond)
	msr.SetProcessingDuration(3, 10*time.Millisecond)

	opts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measurements, err := msr.Play(discipline)
	require.NoError(t, err)

	createGraphs(
		t,
		"Rate divider, even time processing",
		"rate_even",
		opts.HandlersQuantity,
		unbufferedInput,
		measurements,
	)
}

func TestGraphRateUneven(t *testing.T) {
	testGraphRateUneven(t, 1, false)
	testGraphRateUneven(t, 10, false)
	testGraphRateUneven(t, 100, false)
}

func TestGraphRateUnevenUnbuffered(t *testing.T) {
	testGraphRateUneven(t, 1, true)
	testGraphRateUneven(t, 10, true)
	testGraphRateUneven(t, 100, true)
}

func testGraphRateUneven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv(env.EnableGraphs) == "" {
		t.SkipNow()
	}

	msr, err := measuring.NewMeasurer(6*factor, inputCapacity(unbufferedInput)...)
	require.NoError(t, err)

	msr.AddWrite(1, 390*factor)

	msr.AddWrite(2, 250*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 100*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 4*time.Second)
	msr.AddWrite(2, 150*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 300*factor)

	msr.AddWrite(3, 1000*factor)
	msr.AddWaitDevastation(3)
	msr.AddDelay(3, 7*time.Second)
	msr.AddWrite(3, 3500*factor)

	msr.SetProcessingDuration(1, 100*time.Millisecond)
	msr.SetProcessingDuration(2, 50*time.Millisecond)
	msr.SetProcessingDuration(3, 10*time.Millisecond)

	opts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measurements, err := msr.Play(discipline)
	require.NoError(t, err)

	createGraphs(
		t,
		"Rate divider, uneven time processing",
		"rate_uneven",
		opts.HandlersQuantity,
		unbufferedInput,
		measurements,
	)
}

func TestGraphRate721Even(t *testing.T) {
	testGraphRate721Even(t, 1, false)
	testGraphRate721Even(t, 10, false)
	testGraphRate721Even(t, 100, false)
}

func TestGraphRate721EvenUnbuffered(t *testing.T) {
	testGraphRate721Even(t, 1, true)
	testGraphRate721Even(t, 10, true)
	testGraphRate721Even(t, 100, true)
}

func testGraphRate721Even(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv(env.EnableGraphs) == "" {
		t.SkipNow()
	}

	msr, err := measuring.NewMeasurer(10*factor, inputCapacity(unbufferedInput)...)
	require.NoError(t, err)

	msr.AddWrite(1, 5200*factor)

	msr.AddWrite(2, 1500*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 750*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 4*time.Second)
	msr.AddWrite(2, 700*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 2000*factor)

	msr.AddWrite(7, 1000*factor)
	msr.AddWaitDevastation(7)
	msr.AddDelay(7, 7*time.Second)
	msr.AddWrite(7, 7000*factor)

	msr.SetProcessingDuration(1, 10*time.Millisecond)
	msr.SetProcessingDuration(2, 10*time.Millisecond)
	msr.SetProcessingDuration(7, 10*time.Millisecond)

	opts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measurements, err := msr.Play(discipline)
	require.NoError(t, err)

	createGraphs(
		t,
		"Rate divider, even time processing",
		"rate_721_even",
		opts.HandlersQuantity,
		unbufferedInput,
		measurements,
	)
}

func TestGraphRate721Uneven(t *testing.T) {
	testGraphRate721Uneven(t, 1, false)
	testGraphRate721Uneven(t, 10, false)
	testGraphRate721Uneven(t, 100, false)
}

func TestGraphRate721UnevenUnbuffered(t *testing.T) {
	testGraphRate721Uneven(t, 1, true)
	testGraphRate721Uneven(t, 10, true)
	testGraphRate721Uneven(t, 100, true)
}

func testGraphRate721Uneven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv(env.EnableGraphs) == "" {
		t.SkipNow()
	}

	msr, err := measuring.NewMeasurer(10*factor, inputCapacity(unbufferedInput)...)
	require.NoError(t, err)

	msr.AddWrite(1, 500*factor)

	msr.AddWrite(2, 250*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 100*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 4*time.Second)
	msr.AddWrite(2, 130*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 500*factor)

	msr.AddWrite(7, 1000*factor)
	msr.AddWaitDevastation(7)
	msr.AddDelay(7, 6*time.Second)
	msr.AddWrite(7, 6800*factor)

	msr.SetProcessingDuration(1, 100*time.Millisecond)
	msr.SetProcessingDuration(2, 50*time.Millisecond)
	msr.SetProcessingDuration(7, 10*time.Millisecond)

	opts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measurements, err := msr.Play(discipline)
	require.NoError(t, err)

	createGraphs(
		t,
		"Rate divider, uneven time processing",
		"rate_721_uneven",
		opts.HandlersQuantity,
		unbufferedInput,
		measurements,
	)
}

func TestGraphUnmanagedEven(t *testing.T) {
	testGraphUnmanagedEven(t, 1, false)
	testGraphUnmanagedEven(t, 10, false)
	testGraphUnmanagedEven(t, 100, false)
}

func TestGraphUnmanagedEvenUnbuffered(t *testing.T) {
	testGraphUnmanagedEven(t, 1, true)
	testGraphUnmanagedEven(t, 10, true)
	testGraphUnmanagedEven(t, 100, true)
}

func testGraphUnmanagedEven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv(env.EnableGraphs) == "" {
		t.SkipNow()
	}

	msr, err := measuring.NewMeasurer(6*factor, inputCapacity(unbufferedInput)...)
	require.NoError(t, err)

	msr.AddWrite(1, 4000*factor)

	msr.AddWrite(2, 500*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 500*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 4*time.Second)
	msr.AddWrite(2, 1000*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 2000*factor)

	msr.AddWrite(3, 500*factor)
	msr.AddWaitDevastation(3)
	msr.AddDelay(3, 5*time.Second)
	msr.AddWrite(3, 4000*factor)

	msr.SetProcessingDuration(1, 10*time.Millisecond)
	msr.SetProcessingDuration(2, 10*time.Millisecond)
	msr.SetProcessingDuration(3, 10*time.Millisecond)

	opts := unmanaged.Opts[uint]{
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := unmanaged.New(opts)
	require.NoError(t, err)

	measurements, err := msr.Play(discipline)
	require.NoError(t, err)

	createGraphs(
		t,
		"Unmanaged, even time processing",
		"unmanaged_even",
		opts.HandlersQuantity,
		unbufferedInput,
		measurements,
	)
}

func TestGraphUnmanagedUneven(t *testing.T) {
	testGraphUnmanagedUneven(t, 1, false)
	testGraphUnmanagedUneven(t, 10, false)
	testGraphUnmanagedUneven(t, 100, false)
}

func TestGraphUnmanagedUnevenUnbuffered(t *testing.T) {
	testGraphUnmanagedUneven(t, 1, true)
	testGraphUnmanagedUneven(t, 10, true)
	testGraphUnmanagedUneven(t, 100, true)
}

func testGraphUnmanagedUneven(t *testing.T, factor uint, unbufferedInput bool) {
	if os.Getenv(env.EnableGraphs) == "" {
		t.SkipNow()
	}

	msr, err := measuring.NewMeasurer(6*factor, inputCapacity(unbufferedInput)...)
	require.NoError(t, err)

	msr.AddWrite(1, 500*factor)

	msr.AddWrite(2, 100*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 100*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 200*factor)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 400*factor)

	msr.AddWrite(3, 100*factor)
	msr.AddWaitDevastation(3)
	msr.AddDelay(3, 6*time.Second)
	msr.AddWrite(3, 1350*factor)

	msr.SetProcessingDuration(1, 100*time.Millisecond)
	msr.SetProcessingDuration(2, 50*time.Millisecond)
	msr.SetProcessingDuration(3, 10*time.Millisecond)

	opts := unmanaged.Opts[uint]{
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := unmanaged.New(opts)
	require.NoError(t, err)

	measurements, err := msr.Play(discipline)
	require.NoError(t, err)

	createGraphs(
		t,
		"Unmanaged, uneven time processing",
		"unmanaged_uneven",
		opts.HandlersQuantity,
		unbufferedInput,
		measurements,
	)
}

func TestGraphDividingError(t *testing.T) {
	testGraphDividingError(t, 6)
	testGraphDividingError(t, 7)
	testGraphDividingError(t, 8)
	testGraphDividingError(t, 9)
	testGraphDividingError(t, 10)
	testGraphDividingError(t, 11)
	testGraphDividingError(t, 12)
}

func testGraphDividingError(t *testing.T, handlersQuantity uint) {
	if os.Getenv(env.EnableGraphs) == "" {
		t.SkipNow()
	}

	msr, err := measuring.NewMeasurer(handlersQuantity)
	require.NoError(t, err)

	msr.AddWrite(1, 4000)

	msr.AddWrite(2, 500)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 500)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 4*time.Second)
	msr.AddWrite(2, 1000)
	msr.AddWaitDevastation(2)
	msr.AddDelay(2, 2*time.Second)
	msr.AddWrite(2, 2000)

	msr.AddWrite(3, 500)
	msr.AddWaitDevastation(3)
	msr.AddDelay(3, 5*time.Second)
	msr.AddWrite(3, 4000)

	msr.AddWrite(4, 500)
	msr.AddWaitDevastation(3)
	msr.AddDelay(4, 5*time.Second)
	msr.AddWrite(4, 4000)

	msr.SetProcessingDuration(1, 10*time.Millisecond)
	msr.SetProcessingDuration(2, 10*time.Millisecond)
	msr.SetProcessingDuration(3, 10*time.Millisecond)
	msr.SetProcessingDuration(4, 10*time.Millisecond)

	opts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measurements, err := msr.Play(discipline)
	require.NoError(t, err)

	createGraphs(
		t,
		"Fair divider, even time processing, significant dividing error",
		"fair_even_dividing_error",
		opts.HandlersQuantity,
		false,
		measurements,
	)
}

func createGraphs(
	t *testing.T,
	subtitleBase string,
	filePrefix string,
	handlersQuantity uint,
	unbufferedInput bool,
	measurements []measuring.Measure,
) {
	const (
		graphInterval = 100 * time.Millisecond
		graphTimeUnit = time.Second
	)

	received := slices.DeleteFunc(slices.Clone(measurements), measuring.KeepReceived)

	dqot, dqotX := research.QotToLineChart(
		research.QuantityPerInterval(received, graphInterval),
		graphTimeUnit,
	)

	ipot, ipotX := research.QotToLineChart(
		research.InProcessing(measurements, graphInterval),
		graphTimeUnit,
	)

	subtitle := fmt.Sprintf(
		"%s, handlers quantity: %d, unbuffered: %t, time: %s",
		subtitleBase,
		handlersQuantity,
		unbufferedInput,
		time.Now().Format(time.RFC3339),
	)

	baseName := fmt.Sprintf(
		"graph_%s_%d_unbuffered_%t",
		filePrefix,
		handlersQuantity,
		unbufferedInput,
	)

	createLineGraph(
		t,
		"Graph of data receiving",
		subtitle,
		baseName+"_data_receiving.html",
		dqot,
		dqotX,
	)

	createLineGraph(
		t,
		"Graph of data being in processing",
		subtitle,
		baseName+"_in_processing.html",
		ipot,
		ipotX,
	)
}

func createLineGraph(
	t *testing.T,
	title string,
	subtitle string,
	fileName string,
	serieses map[uint][]chartsopts.LineData,
	abscissa []int,
) {
	if len(serieses) == 0 {
		return
	}

	chart := charts.NewLine()

	chart.SetGlobalOptions(
		charts.WithTitleOpts(
			chartsopts.Title{
				Title:    title,
				Subtitle: subtitle,
			},
		),
	)

	addLineSeries(chart.SetXAxis(abscissa), serieses)

	file, err := os.Create(fileName)
	require.NoError(t, err)

	defer file.Close()

	err = chart.Render(file)
	require.NoError(t, err)
}

func addLineSeries(line *charts.Line, serieses map[uint][]chartsopts.LineData) {
	priorities := make([]uint, 0, len(serieses))

	for priority := range serieses {
		priorities = append(priorities, priority)
	}

	slices.SortFunc(priorities, Compare)

	for _, priority := range priorities {
		line.AddSeries(
			strconv.FormatUint(uint64(priority), consts.DecimalBase),
			serieses[priority],
		)
	}
}

func inputCapacity(unbufferedInput bool) []uint {
	if unbufferedInput {
		return []uint{0}
	}

	return nil
}
