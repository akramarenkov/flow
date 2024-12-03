package measuring

import (
	"maps"
	"slices"

	"github.com/akramarenkov/span"
)

// Validates measurements.
func isCorrectMeasurements(measurements []Measure, expectedSpans map[uint]span.Span[uint]) error {
	if len(measurements) < measurementsPerDataItem {
		return ErrMeasurementsIsIncomplete
	}

	slices.SortFunc(measurements, CompareItem)

	first := 0
	last := len(measurements) - 1

	kinds := make(map[Kind]struct{}, measurementsPerDataItem)
	kinds[measurements[first].Kind] = struct{}{}

	spans := make(map[uint]span.Span[uint])

	spans[measurements[first].Priority] = span.Span[uint]{
		Begin: measurements[first].Item,
	}

	spans[measurements[last].Priority] = span.Span[uint]{
		End: measurements[last].Item,
	}

	for id, current := range measurements[:last] {
		next := measurements[id+1]
		diff := next.Item - current.Item

		if diff == 0 {
			if _, exists := kinds[next.Kind]; exists {
				return ErrMeasureDuplicated
			}

			if next.Priority != current.Priority {
				return ErrDataMixedBetweenPriorities
			}

			kinds[next.Kind] = struct{}{}

			continue
		}

		if diff != 1 {
			return ErrDataPartiallyLost
		}

		if id+1 == last {
			return ErrMeasurementsIsIncomplete
		}

		if _, exists := kinds[KindCompleted]; !exists {
			return ErrMeasurementsIsIncomplete
		}

		if _, exists := kinds[KindProcessed]; !exists {
			return ErrMeasurementsIsIncomplete
		}

		if _, exists := kinds[KindReceived]; !exists {
			return ErrMeasurementsIsIncomplete
		}

		if len(kinds) != measurementsPerDataItem {
			return ErrUnexpectedMeasureKind
		}

		clear(kinds)

		kinds[next.Kind] = struct{}{}

		if next.Priority != current.Priority {
			spans[current.Priority] = span.Span[uint]{
				Begin: spans[current.Priority].Begin,
				End:   current.Item,
			}

			spans[next.Priority] = span.Span[uint]{
				Begin: next.Item,
				End:   spans[next.Priority].End,
			}
		}
	}

	if !maps.Equal(expectedSpans, spans) {
		return ErrDataMixedBetweenPriorities
	}

	return nil
}
