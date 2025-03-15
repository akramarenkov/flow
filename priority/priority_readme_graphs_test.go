package priority

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/akramarenkov/flow/internal/env"
	"github.com/akramarenkov/flow/priority/divider"
	"github.com/akramarenkov/flow/priority/internal/measuring"
	"github.com/akramarenkov/flow/priority/internal/research"
	"github.com/akramarenkov/flow/priority/internal/unmanaged"

	"github.com/stretchr/testify/require"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func TestReadmeGraph(t *testing.T) {
	t.Run(
		"equaling",
		func(t *testing.T) {
			t.Parallel()
			testReadmeGraph(t, true)
		},
	)

	t.Run(
		"unmanaged",
		func(t *testing.T) {
			t.Parallel()
			testReadmeGraph(t, false)
		},
	)
}

func testReadmeGraph(t *testing.T, equaling bool) {
	if os.Getenv(env.EnableGraphs) == "" {
		t.SkipNow()
	}

	msr, err := measuring.NewMeasurer(6)
	require.NoError(t, err)

	msr.AddWrite(1, 500)
	msr.AddWrite(2, 500)
	msr.AddWrite(3, 500)

	msr.SetProcessingDuration(1, 100*time.Millisecond)
	msr.SetProcessingDuration(2, 50*time.Millisecond)
	msr.SetProcessingDuration(3, 10*time.Millisecond)

	graphInterval := 100 * time.Millisecond
	graphTimeUnit := time.Second
	graphTimeUnitName := "seconds"

	if equaling {
		opts := Opts[uint]{
			Divider:          divider.Fair,
			HandlersQuantity: msr.HandlersQuantity(),
			Inputs:           msr.Inputs(),
		}

		discipline, err := New(opts)
		require.NoError(t, err)

		measurements, err := msr.Play(discipline)
		require.NoError(t, err)

		createReadmeGraph(
			t,
			"doc/different-processing-time-equaling.svg",
			measurements,
			graphInterval,
			graphTimeUnit,
			graphTimeUnitName,
		)

		return
	}

	opts := unmanaged.Opts[uint]{
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := unmanaged.New(opts)
	require.NoError(t, err)

	measurements, err := msr.Play(discipline)
	require.NoError(t, err)

	createReadmeGraph(
		t,
		"doc/different-processing-time-unmanaged.svg",
		measurements,
		graphInterval,
		graphTimeUnit,
		graphTimeUnitName,
	)
}

func createReadmeGraph(
	t *testing.T,
	fileName string,
	measurements []measuring.Measure,
	graphInterval time.Duration,
	graphTimeUnit time.Duration,
	graphTimeUnitName string,
) {
	received := slices.DeleteFunc(slices.Clone(measurements), measuring.KeepReceived)
	researched := research.QuantityPerInterval(received, graphInterval)

	series := make([]chart.Series, 0, len(researched))

	// To keep the legends in the same order
	priorities := slices.SortedFunc(maps.Keys(researched), Compare)

	for _, priority := range priorities {
		xaxis := make([]float64, len(researched[priority]))
		yaxis := make([]float64, len(researched[priority]))

		for id, item := range researched[priority] {
			xaxis[id] = float64(item.Time) / float64(graphTimeUnit)
			yaxis[id] = float64(item.Quantity)
		}

		srs := chart.ContinuousSeries{
			Name:    fmt.Sprintf("Data of priority %d", priority),
			XValues: xaxis,
			YValues: yaxis,
			Style:   chart.Style{StrokeWidth: 4},
		}

		series = append(series, srs)
	}

	graph := chart.Chart{
		Title:        "Graph of data receiving",
		ColorPalette: readmeColorPalette{},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  50,
				Left: 140,
			},
			FillColor: chart.ColorTransparent,
		},
		Canvas: chart.Style{
			FillColor: chart.ColorTransparent,
		},
		XAxis: chart.XAxis{
			Name: "Time, " + graphTimeUnitName,
		},
		YAxis: chart.YAxis{
			Name: "Quantity of data items received by handlers, pieces",
		},
		Series: series,
	}

	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph, chart.Style{FillColor: chart.ColorTransparent}),
	}

	file, err := os.Create(fileName)
	require.NoError(t, err)

	defer file.Close()

	err = graph.Render(chart.SVG, file)
	require.NoError(t, err)
}

type readmeColorPalette struct{}

func (readmeColorPalette) BackgroundColor() drawing.Color {
	return chart.DefaultColorPalette.BackgroundColor()
}

func (readmeColorPalette) BackgroundStrokeColor() drawing.Color {
	return chart.DefaultColorPalette.BackgroundStrokeColor()
}

func (readmeColorPalette) CanvasColor() drawing.Color {
	return chart.DefaultColorPalette.CanvasColor()
}

func (readmeColorPalette) CanvasStrokeColor() drawing.Color {
	return chart.DefaultColorPalette.CanvasStrokeColor()
}

func (readmeColorPalette) AxisStrokeColor() drawing.Color {
	return chart.DefaultColorPalette.AxisStrokeColor()
}

func (readmeColorPalette) TextColor() drawing.Color {
	return chart.DefaultColorPalette.TextColor()
}

func (readmeColorPalette) GetSeriesColor(index int) drawing.Color {
	colors := []drawing.Color{
		{R: 0x54, G: 0x70, B: 0xc6, A: 255},
		{R: 0x91, G: 0xcc, B: 0x75, A: 255},
		{R: 0xfa, G: 0xc8, B: 0x58, A: 255},
		{R: 0xee, G: 0x66, B: 0x66, A: 255},
		{R: 0x73, G: 0xc0, B: 0xde, A: 255},
		{R: 0x3b, G: 0xa2, B: 0x72, A: 255},
		{R: 0xfc, G: 0x84, B: 0x52, A: 255},
		{R: 0x9a, G: 0x60, B: 0xb4, A: 255},
		{R: 0xea, G: 0x7c, B: 0xcc, A: 255},
	}

	return colors[index%len(colors)]
}
