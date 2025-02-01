package join_test

import (
	"fmt"
	"time"

	"github.com/akramarenkov/flow/join"
)

func ExampleDiscipline() {
	data := []int{
		1, 2, 3, 4, 5, 6, 7, 8,
		9, 10, 11, 12, 13, 14, 15, 16,
		17, 18, 19, 20, 21, 22, 23, 24,
		25, 26, 27,
	}

	// Preferably input channel should be buffered for performance reasons.
	// Optimal capacity is in the range of 1 to 3 size of join
	input := make(chan int, 10)

	opts := join.Opts[int]{
		Input:    input,
		JoinSize: 10,
		Timeout:  time.Second,
	}

	discipline, err := join.New(opts)
	if err != nil {
		panic(err)
	}

	go func() {
		defer close(input)

		for _, item := range data {
			input <- item
		}
	}()

	for joined := range discipline.Output() {
		fmt.Println(joined)
		discipline.Release()
	}
	// Output:
	// [1 2 3 4 5 6 7 8 9 10]
	// [11 12 13 14 15 16 17 18 19 20]
	// [21 22 23 24 25 26 27]
}
