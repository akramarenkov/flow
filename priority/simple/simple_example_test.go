package simple_test

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/akramarenkov/flow/priority"
	"github.com/akramarenkov/flow/priority/divider"
	"github.com/akramarenkov/flow/priority/simple"
	"github.com/akramarenkov/flow/priority/types"

	"github.com/guptarohit/asciigraph"
)

func ExampleDiscipline() {
	handlersQuantity := uint(100)
	itemsQuantity := 10000
	// Preferably input channels should be buffered for performance reasons.
	// Optimal capacity is in the range of 1 to 3 times of quantity of data handlers
	inputCapacity := handlersQuantity
	processingDuration := 10 * time.Millisecond
	graphInterval := 100 * time.Millisecond
	graphRange := 5 * time.Second

	inputs := map[uint]chan int{
		70: make(chan int, inputCapacity),
		20: make(chan int, inputCapacity),
		10: make(chan int, inputCapacity),
	}

	// Used only in this example for measuring receiving of data items
	type measure struct {
		Priority uint
		Time     time.Duration
	}

	compareTime := func(first, second measure) int {
		return cmp.Compare(first.Time, second.Time)
	}

	// Channel capacity is equal to the total quantity of input data in order to
	// minimize delays in collecting measurements
	measurements := make(chan measure, itemsQuantity*len(inputs))

	startedAt := time.Now()

	handle := func(prioritized types.Prioritized[int]) {
		// Data item processing
		measurement := measure{
			Priority: prioritized.Priority,
			Time:     time.Since(startedAt),
		}

		time.Sleep(processingDuration)

		measurements <- measurement
	}

	// For equaling use divider.Fair divider, for prioritization use divider.Rate
	// divider or custom divider
	opts := simple.Opts[int]{
		Divider:          divider.Rate,
		Handle:           handle,
		HandlersQuantity: handlersQuantity,
	}

	for priority, channel := range inputs {
		if err := opts.AddInput(priority, channel); err != nil {
			panic(err)
		}
	}

	discipline, err := simple.New(opts)
	if err != nil {
		panic(err)
	}

	// Running writers, that write data items to input channels
	for _, input := range inputs {
		go func() {
			defer close(input)

			for item := range itemsQuantity {
				input <- item
			}
		}()
	}

	// Waiting for completion of the discipline, and also writers and handlers
	if err := <-discipline.Err(); err != nil {
		fmt.Println("An error was received: ", err)
	}

	close(measurements)

	received := make(map[uint][]measure, len(inputs))

	// Receiving measurements
	for item := range measurements {
		received[item.Priority] = append(received[item.Priority], item)
	}

	// Sorting measurements by time for further research
	for _, measurements := range received {
		slices.SortFunc(measurements, compareTime)
	}

	// Calculating quantity of data items received by handlers over time
	quantities := make(map[uint][]float64)

	for span := time.Duration(0); span <= graphRange; span += graphInterval {
		for priority, measurements := range received {
			quantity := float64(0)

			for _, measure := range measurements {
				if measure.Time < span-graphInterval {
					continue
				}

				if measure.Time >= span {
					break
				}

				quantity++
			}

			quantities[priority] = append(quantities[priority], quantity)
		}
	}

	// Preparing research data for plot
	serieses := make([][]float64, 0, len(quantities))
	priorities := make([]uint, 0, len(quantities))
	legends := make([]string, 0, len(quantities))

	for priority := range quantities {
		priorities = append(priorities, priority)
	}

	// To keep the legends in the same order
	slices.SortFunc(priorities, priority.Compare)

	for _, priority := range priorities {
		serieses = append(serieses, quantities[priority])
		legends = append(legends, strconv.FormatUint(uint64(priority), 10))
	}

	graph := asciigraph.PlotMany(
		serieses,
		asciigraph.Height(10),
		asciigraph.Caption("Quantity of data items received by handlers over time"),
		asciigraph.SeriesColors(asciigraph.Red, asciigraph.Green, asciigraph.Blue),
		asciigraph.SeriesLegends(legends...),
	)

	fmt.Fprintln(os.Stderr, graph)
	fmt.Println("See graph")
	// Output:
	// See graph
}
