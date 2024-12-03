package measuring

import (
	"math"
	"testing"

	"github.com/akramarenkov/span"
	"github.com/stretchr/testify/require"
)

func TestIsCorrectMeasurements(t *testing.T) {
	spans := map[uint]span.Span[uint]{
		3: {Begin: 0, End: 0},
		2: {Begin: 1, End: 1},
	}

	measurements := []Measure{
		{Item: 0, Kind: KindReceived, Priority: 3},
		{Item: 1, Kind: KindProcessed, Priority: 2},
		{Item: 0, Kind: KindCompleted, Priority: 3},
		{Item: 1, Kind: KindReceived, Priority: 2},
		{Item: 0, Kind: KindProcessed, Priority: 3},
		{Item: 1, Kind: KindCompleted, Priority: 2},
	}

	require.NoError(t, isCorrectMeasurements(measurements, spans))
}

func TestIsCorrectMeasurementsIncomplete(t *testing.T) {
	spans := map[uint]span.Span[uint]{
		3: {Begin: 0, End: 0},
		2: {Begin: 1, End: 1},
	}

	require.Error(t, isCorrectMeasurements(nil, spans))
	require.Error(t, isCorrectMeasurements([]Measure{}, spans))
	require.Error(t, isCorrectMeasurements([]Measure{{}}, spans))
	require.Error(t, isCorrectMeasurements([]Measure{{}, {}}, spans))

	measurements := []Measure{
		{Item: 0, Kind: KindReceived, Priority: 3},
		{Item: 1, Kind: KindProcessed, Priority: 2},
		{Item: 0, Kind: KindCompleted, Priority: 3},
		{Item: 2, Kind: KindReceived, Priority: 2},
		{Item: 0, Kind: KindProcessed, Priority: 3},
		{Item: 1, Kind: KindCompleted, Priority: 2},
	}

	require.Error(t, isCorrectMeasurements(measurements, spans))

	measurements = []Measure{
		{Item: 0, Kind: KindReceived, Priority: 3},
		{Item: 1, Kind: KindProcessed, Priority: 2},
		{Item: 1, Kind: KindReceived, Priority: 2},
		{Item: 0, Kind: KindProcessed, Priority: 3},
		{Item: 1, Kind: KindCompleted, Priority: 2},
	}

	require.Error(t, isCorrectMeasurements(measurements, spans))

	measurements = []Measure{
		{Item: 0, Kind: KindReceived, Priority: 3},
		{Item: 1, Kind: KindProcessed, Priority: 2},
		{Item: 0, Kind: KindCompleted, Priority: 3},
		{Item: 1, Kind: KindReceived, Priority: 2},
		{Item: 1, Kind: KindCompleted, Priority: 2},
	}

	require.Error(t, isCorrectMeasurements(measurements, spans))

	measurements = []Measure{
		{Item: 1, Kind: KindProcessed, Priority: 2},
		{Item: 0, Kind: KindCompleted, Priority: 3},
		{Item: 1, Kind: KindReceived, Priority: 2},
		{Item: 0, Kind: KindProcessed, Priority: 3},
		{Item: 1, Kind: KindCompleted, Priority: 2},
	}

	require.Error(t, isCorrectMeasurements(measurements, spans))
}

func TestIsCorrectMeasurementsDuplicated(t *testing.T) {
	spans := map[uint]span.Span[uint]{
		3: {Begin: 0, End: 0},
		2: {Begin: 1, End: 1},
	}

	measurements := []Measure{
		{Item: 0, Kind: KindReceived, Priority: 3},
		{Item: 1, Kind: KindProcessed, Priority: 2},
		{Item: 0, Kind: KindCompleted, Priority: 3},
		{Item: 1, Kind: KindReceived, Priority: 2},
		{Item: 0, Kind: KindReceived, Priority: 3},
		{Item: 1, Kind: KindCompleted, Priority: 2},
	}

	require.Error(t, isCorrectMeasurements(measurements, spans))
}

func TestIsCorrectMeasurementsMixed(t *testing.T) {
	spans := map[uint]span.Span[uint]{
		3: {Begin: 0, End: 0},
		2: {Begin: 1, End: 1},
	}

	measurements := []Measure{
		{Item: 0, Kind: KindReceived, Priority: 3},
		{Item: 1, Kind: KindProcessed, Priority: 2},
		{Item: 0, Kind: KindCompleted, Priority: 3},
		{Item: 1, Kind: KindReceived, Priority: 2},
		{Item: 0, Kind: KindProcessed, Priority: 2},
		{Item: 1, Kind: KindCompleted, Priority: 2},
	}

	require.Error(t, isCorrectMeasurements(measurements, spans))

	measurements = []Measure{
		{Item: 1, Kind: KindReceived, Priority: 3},
		{Item: 0, Kind: KindProcessed, Priority: 2},
		{Item: 1, Kind: KindCompleted, Priority: 3},
		{Item: 0, Kind: KindReceived, Priority: 2},
		{Item: 1, Kind: KindProcessed, Priority: 3},
		{Item: 0, Kind: KindCompleted, Priority: 2},
	}

	require.Error(t, isCorrectMeasurements(measurements, spans))
}

func TestIsCorrectMeasurementsLost(t *testing.T) {
	spans := map[uint]span.Span[uint]{
		3: {Begin: 0, End: 0},
		2: {Begin: 1, End: 1},
	}

	measurements := []Measure{
		{Item: 0, Kind: KindReceived, Priority: 3},
		{Item: 2, Kind: KindProcessed, Priority: 2},
		{Item: 0, Kind: KindCompleted, Priority: 3},
		{Item: 2, Kind: KindReceived, Priority: 2},
		{Item: 0, Kind: KindProcessed, Priority: 3},
		{Item: 2, Kind: KindCompleted, Priority: 2},
	}

	require.Error(t, isCorrectMeasurements(measurements, spans))
}

func TestIsCorrectMeasurementsLostUnexpected(t *testing.T) {
	spans := map[uint]span.Span[uint]{
		3: {Begin: 0, End: 0},
		2: {Begin: 1, End: 1},
	}

	measurements := []Measure{
		{Item: 0, Kind: math.MaxInt, Priority: 3},
		{Item: 0, Kind: KindReceived, Priority: 3},
		{Item: 1, Kind: KindProcessed, Priority: 2},
		{Item: 0, Kind: KindCompleted, Priority: 3},
		{Item: 1, Kind: KindReceived, Priority: 2},
		{Item: 0, Kind: KindProcessed, Priority: 3},
		{Item: 1, Kind: KindCompleted, Priority: 2},
	}

	require.Error(t, isCorrectMeasurements(measurements, spans))
}

func BenchmarkIsCorrectMeasurements(b *testing.B) {
	spans := map[uint]span.Span[uint]{
		3: {Begin: 0, End: 0},
		2: {Begin: 1, End: 1},
	}

	measurements := []Measure{
		{Item: 0, Kind: KindReceived, Priority: 3},
		{Item: 1, Kind: KindProcessed, Priority: 2},
		{Item: 0, Kind: KindCompleted, Priority: 3},
		{Item: 1, Kind: KindReceived, Priority: 2},
		{Item: 0, Kind: KindProcessed, Priority: 3},
		{Item: 1, Kind: KindCompleted, Priority: 2},
	}

	var err error

	for range b.N {
		err = isCorrectMeasurements(measurements, spans)
	}

	require.NoError(b, err)
}
