package utils

import (
	"math"
	"slices"

	"github.com/akramarenkov/flow/internal/consts"
	"github.com/akramarenkov/flow/priority/divider"
	"github.com/akramarenkov/flow/priority/internal/distrib"
	"github.com/akramarenkov/flow/priority/internal/priocombs"
	"github.com/akramarenkov/flow/priority/priolist"
	"github.com/akramarenkov/safe"
)

func isNonFatalConfig(
	priorities []uint,
	divider divider.Divider,
	handlersQuantity uint,
) (bool, error) {
	if err := priolist.IsValid(priorities); err != nil {
		return false, err
	}

	for combination := range priocombs.Iter(priorities) {
		distribution := make(map[uint]uint)

		divider(combination, handlersQuantity, distribution)

		if !distrib.IsFilled(distribution) {
			return false, nil
		}
	}

	return true, nil
}

// Due to the imperfection of the dividing function and working with integers (since
// the quantity of data handlers is an integer), large errors can occur when
// distributing handlers by priorities, especially for small quantity of data handlers.
// This function allows you to determine that with the specified combination of
// priorities, the dividing function and the quantity of data handlers, the distribution
// error does not cause stop processing of one or more priorities (for none of the
// priorities, the distributed quantity of data handlers is not equal to zero).
//
// This function panics if the priority list has zero length or contains zero priority.
func IsNonFatalConfig(priorities []uint, divider divider.Divider, handlersQuantity uint) bool {
	priorities = slices.SortedFunc(slices.Values(priorities), priolist.Compare)

	conclusion, err := isNonFatalConfig(priorities, divider, handlersQuantity)
	if err != nil {
		panic(err)
	}

	return conclusion
}

// Picks up the minimum quantity of data handlers for which the distribution error does
// not cause stop processing of one or more priorities.
//
// If it is impossible to pickup the quantity of data handlers, zero is returned.
//
// This function panics if the priority list has zero length or contains zero priority.
func PickUpMinNonFatalQuantity(priorities []uint, divider divider.Divider, maxHandlersQuantity uint) uint {
	priorities = slices.SortedFunc(slices.Values(priorities), priolist.Compare)

	for quantity := range safe.Inc(1, maxHandlersQuantity) {
		nonFatal, err := isNonFatalConfig(priorities, divider, quantity)
		if err != nil {
			panic(err)
		}

		if nonFatal {
			return quantity
		}
	}

	return 0
}

// Picks up the maximum quantity of data handlers for which the distribution error does
// not cause stop processing of one or more priorities.
//
// If it is impossible to pickup the quantity of data handlers, zero is returned.
//
// This function panics if the priority list has zero length or contains zero priority.
func PickUpMaxNonFatalQuantity(priorities []uint, divider divider.Divider, maxHandlersQuantity uint) uint {
	priorities = slices.SortedFunc(slices.Values(priorities), priolist.Compare)

	for quantity := range safe.Dec(maxHandlersQuantity, 1) {
		nonFatal, err := isNonFatalConfig(priorities, divider, quantity)
		if err != nil {
			panic(err)
		}

		if nonFatal {
			return quantity
		}
	}

	return 0
}

// Acceptable relative difference is specified as a percentage.
func isSuitableDiff(
	reference map[uint]uint,
	distribution map[uint]uint,
	referenceRatio uint,
	acceptableRelativeDiff float64,
) bool {
	for priority, referenceQuantity := range reference {
		quantity := distribution[priority]

		diff := 1.0 - float64(referenceRatio)*float64(quantity)/float64(referenceQuantity)
		diff = consts.HundredPercent * math.Abs(diff)

		if diff > acceptableRelativeDiff {
			return false
		}
	}

	return true
}

func isSuitableConfig(
	priorities []uint,
	divider divider.Divider,
	handlersQuantity uint,
	acceptableRelativeDiff float64,
) (bool, error) {
	const referenceRatioFactor = 100

	if err := priolist.IsValid(priorities); err != nil {
		return false, err
	}

	// Increases the chances of getting correct reference values
	prioritiesSum, err := safe.AddMU(priorities...)
	if err != nil {
		return false, err
	}

	referenceRatio, err := safe.Mul(prioritiesSum, referenceRatioFactor)
	if err != nil {
		return false, err
	}

	referenceHandlersQuantity, err := safe.Mul(handlersQuantity, referenceRatio)
	if err != nil {
		return false, err
	}

	for combination := range priocombs.Iter(priorities) {
		reference := make(map[uint]uint)
		distribution := make(map[uint]uint)

		divider(combination, referenceHandlersQuantity, reference)

		// It is assumed that there is a bug in the division function, in which it
		// either always, or for some combination of priorities, returns 0
		if !distrib.IsFilled(reference) {
			return false, nil
		}

		divider(combination, handlersQuantity, distribution)

		suitable := isSuitableDiff(
			reference,
			distribution,
			referenceRatio,
			acceptableRelativeDiff,
		)

		if !suitable {
			return false, nil
		}
	}

	return true, nil
}

// Due to the imperfection of the dividing function and working with integers (since
// the quantity of data handlers is an integer), large errors can occur when distributing
// handlers by priority, especially for small quantity of data handlers. This function allows
// you to determine that with the specified combination of priorities, the dividing
// function and the quantity of data handlers, the distribution error does not exceed
// the limit, specified as a percentage.
func IsSuitableConfig(
	priorities []uint,
	divider divider.Divider,
	handlersQuantity uint,
	acceptableRelativeDiff float64,
) bool {
	priorities = slices.SortedFunc(slices.Values(priorities), priolist.Compare)

	conclusion, err := isSuitableConfig(
		priorities,
		divider,
		handlersQuantity,
		acceptableRelativeDiff,
	)
	if err != nil {
		panic(err)
	}

	return conclusion
}

// Picks up the minimum quantity of data handlers for which the division error does not
// exceed the limit, specified as a percentage.
//
// If it is impossible to pickup the quantity of data handlers, zero is returned.
func PickUpMinSuitableQuantity(
	priorities []uint,
	divider divider.Divider,
	maxHandlersQuantity uint,
	acceptableRelativeDiff float64,
) uint {
	priorities = slices.SortedFunc(slices.Values(priorities), priolist.Compare)

	for quantity := range safe.Iter(1, maxHandlersQuantity) {
		suitable, err := isSuitableConfig(
			priorities,
			divider,
			quantity,
			acceptableRelativeDiff,
		)
		if err != nil {
			panic(err)
		}

		if suitable {
			return quantity
		}
	}

	return 0
}

// Picks up the maximum quantity of data handlers for which the division error does not
// exceed the limit, specified as a percentage.
//
// If it is impossible to pickup the quantity of data handlers, zero is returned.
func PickUpMaxSuitableQuantity(
	priorities []uint,
	divider divider.Divider,
	maxHandlersQuantity uint,
	acceptableRelativeDiff float64,
) uint {
	priorities = slices.SortedFunc(slices.Values(priorities), priolist.Compare)

	for quantity := range safe.Iter(maxHandlersQuantity, 1) {
		suitable, err := isSuitableConfig(
			priorities,
			divider,
			quantity,
			acceptableRelativeDiff,
		)
		if err != nil {
			panic(err)
		}

		if suitable {
			return quantity
		}
	}

	return 0
}
