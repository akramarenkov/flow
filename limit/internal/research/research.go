// Internal package with research functions that are used for testing.
package research

import (
	"math"
	"slices"
	"strconv"
	"time"

	"github.com/akramarenkov/flow/internal/consts"
	"github.com/akramarenkov/flow/internal/qot"

	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
)

func QuantityPerInterval(
	times []time.Duration,
	intervalsNumber int,
	interval time.Duration,
) ([]qot.QoT, time.Duration) {
	if len(times) == 0 {
		return nil, 0
	}

	slices.Sort(times)

	maxTime := times[len(times)-1]

	if interval == 0 {
		if intervalsNumber == 0 {
			return nil, 0
		}

		// It is necessary that max time falls into the last span,
		// so the interval is rounded up
		interval = time.Duration(math.Ceil(float64(maxTime) / float64(intervalsNumber)))

		maxTimeRecalculated := interval * time.Duration(intervalsNumber)

		// If max time is divided entirely, then add one nanosecond
		// so that it falls into the last span
		if maxTimeRecalculated == maxTime {
			interval += time.Nanosecond
		}
	} else {
		// Number of intervals always turns out to be more by one
		//
		// Due to rounding down during integer division and, if max time is divided
		// entirely, due to the fact that the span takes into account elements
		// strictly smaller than it
		intervalsNumber = int(maxTime/interval) + 1
	}

	quantities := make([]qot.QoT, 0, intervalsNumber)

	edge := 0

	// Interval is added to max span value to be sure that max time falls
	// into the last span
	for span := interval; span <= maxTime+interval; span += interval {
		spanQuantities := uint(0)

		for id, time := range times[edge:] {
			if time >= span {
				edge += id
				break
			}

			spanQuantities++

			// Prevent use of data from the last slice for spans
			// greater than max time + interval
			if id == len(times[edge:])-1 {
				edge += id + 1
			}
		}

		item := qot.QoT{
			Quantity: spanQuantities,
			Time:     span - interval,
		}

		quantities = append(quantities, item)
	}

	// Padding with zero values ​​in case intervals quantity multiplied by
	// interval is greater than max time
	for addition := range intervalsNumber - len(quantities) {
		item := qot.QoT{
			Quantity: 0,
			Time:     maxTime + interval*time.Duration(addition+1),
		}

		quantities = append(quantities, item)
	}

	return quantities, interval
}

func QotToBarChart(quantities []qot.QoT) ([]chartsopts.BarData, []int) {
	serieses := make([]chartsopts.BarData, 0, len(quantities))
	xaxis := make([]int, 0, len(quantities))

	for id, quantity := range quantities {
		item := chartsopts.BarData{
			Name: quantity.Time.String(),
			Tooltip: &chartsopts.Tooltip{
				Show: chartsopts.Bool(true),
			},
			Value: quantity.Quantity,
		}

		serieses = append(serieses, item)
		xaxis = append(xaxis, id)
	}

	return serieses, xaxis
}

func Deviations(times []time.Duration, expected time.Duration) map[int]int {
	const (
		deviationsMin = -100
		deviationsMax = 100
		// from -100% to 100% with 1% step and plus zero
		deviationsLength = deviationsMax - deviationsMin + 1
	)

	if len(times) == 0 {
		return nil
	}

	slices.Sort(times)

	// Used to fix false positive remark by gosec linter
	first := times[0]

	deviations := make(map[int]int, deviationsLength)

	for percent := deviationsMin; percent <= deviationsMax; percent++ {
		deviations[percent] = 0
	}

	calc := func(next time.Duration, current time.Duration) {
		diff := next - current

		deviation := ((diff - expected) * consts.HundredPercent) / expected

		if deviation > consts.HundredPercent {
			deviation = consts.HundredPercent
		}

		if deviation < -consts.HundredPercent {
			deviation = -consts.HundredPercent
		}

		deviations[int(deviation)]++
	}

	calc(first, 0)

	for id := range times {
		if id == len(times)-1 {
			break
		}

		calc(times[id+1], times[id])
	}

	return deviations
}

func DeviationsToBarChart(deviations map[int]int) ([]chartsopts.BarData, []int) {
	serieses := make([]chartsopts.BarData, 0, len(deviations))
	xaxis := make([]int, 0, len(deviations))

	for percent := range deviations {
		xaxis = append(xaxis, percent)
	}

	slices.Sort(xaxis)

	for _, percent := range xaxis {
		item := chartsopts.BarData{
			Name: strconv.Itoa(percent) + "%",
			Tooltip: &chartsopts.Tooltip{
				Show: chartsopts.Bool(true),
			},
			Value: deviations[percent],
		}

		serieses = append(serieses, item)
	}

	return serieses, xaxis
}

func TotalDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	return slices.Max(durations)
}
