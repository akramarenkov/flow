# Limit discipline

## Purpose

Limits the speed of passing data items from the input channel to the output
 channel

The speed limit is set by the **Rate** structure, in which the **Quantity**
 field specifies the number of data items that must pass over the time
 interval specified in the **Interval** field

## Peculiarities

As we know, the speed of 1000 data items per second is, in fact, the same
 speed as 1 data item per millisecond specified in different units of
 measurement

However, the units of measurement affect the distribution of data items
 written to the output channel over time and the performance of the discipline

If the speed is specified as 1000 data items per second, first 1000 data items
 will be written to the output channel, and then a pause will be made equal
 to 1 second minus the time spent writing 1000 data items

If the speed is specified in the form of 1 data item per millisecond, first 1
 data item will be written to the output channel, and then a pause will be
 made equal to 1 millisecond minus the time spent on writing 1 data item

However, the performance of the discipline if the speed is specified in the
 form of 1 data item per millisecond will be lower but the uniformity of the
 distribution of data items over time will be higher

Thus, when choosing units of measurement, you can balance between the uniform
 distribution of data items over time and performance (the maximum achievable
 speed)

Based on measurements, specifying a time interval of less than 10 milliseconds
 significantly reduces performance and accuracy

Under heavy system load, it is not advisable to specify an time interval less
 than 1 second

## Usage

Example:

```go
package main

import (
    "fmt"
    "time"

    "github.com/akramarenkov/flow/limit"
)

func main() {
    data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    // Preferably input channel should be buffered for performance reasons.
    // Optimal capacity is in the range of 1e2 to 1e6
    input := make(chan int, 100)

    opts := limit.Opts[int]{
        Input: input,
        Limit: limit.Rate{
            Interval: time.Second,
            Quantity: 1,
        },
    }

    discipline, err := limit.New(opts)
    if err != nil {
        panic(err)
    }

    output := make([]int, 0, len(data))

    startedAt := time.Now()

    go func() {
        defer close(input)

        for _, item := range data {
            input <- item
        }
    }()

    for item := range discipline.Output() {
        output = append(output, item)
    }

    duration := time.Since(startedAt)
    expected := float64(uint64(len(data))/opts.Limit.Quantity) * float64(opts.Limit.Interval)
    deviation := 0.01

    fmt.Println(duration <= time.Duration(expected*(1.0+deviation)))
    fmt.Println(duration >= time.Duration(expected*(1.0-deviation)))
    fmt.Println(output)
    // Output:
    // true
    // true
    // [1 2 3 4 5 6 7 8 9 10]
}
```
