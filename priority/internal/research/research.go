// Internal package with research functions that are used for testing.
package research

import (
	"slices"
	"time"

	"github.com/akramarenkov/flow/internal/qot"
	"github.com/akramarenkov/flow/priority/internal/measuring"

	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
)

func QuantityPerInterval(
	measurements []measuring.Measure,
	interval time.Duration,
) map[uint][]qot.QoT {
	if len(measurements) == 0 {
		return nil
	}

	slices.SortFunc(measurements, measuring.CompareTime)

	// To see the initial zero values ​​on the graph
	minimum := -interval

	// One interval added to maximum span value to get maximum time into
	// the last span
	//
	// And one more interval added to maximum span value to see the final zero values
	// ​on ​the graph
	expansion := 2
	maximum := measurements[len(measurements)-1].Time + time.Duration(expansion)*interval

	capacity := (maximum - minimum) / interval

	quantities := make(map[uint][]qot.QoT)

	for _, measure := range measurements {
		if _, exists := quantities[measure.Priority]; exists {
			continue
		}

		quantities[measure.Priority] = make([]qot.QoT, 0, capacity)
	}

	edge := 0

	for span := minimum + interval; span <= maximum; span += interval {
		spanQuantities := make(map[uint]uint)

		for id, measure := range measurements[edge:] {
			if measure.Time >= span {
				edge += id
				break
			}

			spanQuantities[measure.Priority]++

			// Prevent use of data from the last slice for spans
			// greater than maximum time + interval
			if id == len(measurements[edge:])-1 {
				edge += id + 1
			}
		}

		for priority, quantity := range spanQuantities {
			item := qot.QoT{
				Quantity: quantity,
				Time:     span - interval,
			}

			quantities[priority] = append(quantities[priority], item)
		}

		for priority := range quantities {
			if _, exists := spanQuantities[priority]; exists {
				continue
			}

			item := qot.QoT{
				Quantity: 0,
				Time:     span - interval,
			}

			quantities[priority] = append(quantities[priority], item)
		}
	}

	return quantities
}

func InProcessing(
	measurements []measuring.Measure,
	interval time.Duration,
) map[uint][]qot.QoT {
	if len(measurements) == 0 {
		return nil
	}

	slices.SortFunc(measurements, measuring.CompareTime)

	// To see the initial zero values ​​on the graph
	minimum := -interval

	// One interval added to maximum span value to get maximum time into
	// the last span
	//
	// And one more interval added to maximum span value to see the final zero values
	// ​on ​the graph
	expansion := 2
	maximum := measurements[len(measurements)-1].Time + time.Duration(expansion)*interval

	capacity := (maximum - minimum) / interval

	quantities := make(map[uint][]qot.QoT)

	for _, measure := range measurements {
		if _, exists := quantities[measure.Priority]; exists {
			continue
		}

		quantities[measure.Priority] = make([]qot.QoT, 0, capacity)
	}

	edge := 0

	receivedQuantities := make(map[uint]map[uint]uint)

	for priority := range quantities {
		receivedQuantities[priority] = make(map[uint]uint)
	}

	for span := minimum + interval; span <= maximum; span += interval {
		for id, measure := range measurements[edge:] {
			if measure.Time >= span {
				edge += id
				break
			}

			switch measure.Kind {
			case measuring.KindReceived:
				receivedQuantities[measure.Priority][measure.Item]++
			case measuring.KindCompleted:
				receivedQuantities[measure.Priority][measure.Item]--
			}

			// Prevent use of data from the last slice for spans
			// greater than maximum time + interval
			if id == len(measurements[edge:])-1 {
				edge += id + 1
			}
		}

		for priority, subset := range receivedQuantities {
			quantity := uint(0)

			for _, amount := range subset {
				quantity += amount
			}

			item := qot.QoT{
				Quantity: quantity,
				Time:     span - interval,
			}

			quantities[priority] = append(quantities[priority], item)
		}
	}

	return quantities
}

func QotToLineChart(
	quantities map[uint][]qot.QoT,
	timeUnit time.Duration,
) (map[uint][]chartsopts.LineData, []int) {
	serieses := make(map[uint][]chartsopts.LineData)
	xaxis := []int(nil)

	for priority := range quantities {
		if _, exists := serieses[priority]; !exists {
			serieses[priority] = make([]chartsopts.LineData, 0, len(quantities[priority]))

			if xaxis == nil {
				xaxis = make([]int, 0, len(quantities[priority]))
			}
		}

		for _, quantity := range quantities[priority] {
			item := chartsopts.LineData{
				Name:  quantity.Time.String(),
				Value: quantity.Quantity,
			}

			serieses[priority] = append(serieses[priority], item)

			if len(xaxis) < len(quantities[priority]) {
				xaxis = append(xaxis, int(quantity.Time/timeUnit))
			}
		}
	}

	return serieses, xaxis
}
