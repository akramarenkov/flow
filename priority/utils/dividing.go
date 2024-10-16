package utils

import (
	"slices"

	"github.com/akramarenkov/flow/priority/divider"
	"github.com/akramarenkov/flow/priority/internal/distrib"
	"github.com/akramarenkov/flow/priority/internal/priocombs"
	"github.com/akramarenkov/flow/priority/priolist"
	"github.com/akramarenkov/safe"
)

func IsCorrectTotalQuantity(priorities []uint, divider divider.Divider, maxQuantity uint) bool {
	priorities = slices.SortedFunc(slices.Values(priorities), priolist.Compare)

	for quantity := range safe.Iter(0, maxQuantity) {
		for combination := range priocombs.Iter(priorities) {
			distribution := make(map[uint]uint)

			divider(combination, quantity, distribution)

			distributed, err := distrib.Quantity(distribution)
			if err != nil {
				return false
			}

			if distributed != quantity {
				return false
			}
		}
	}

	return true
}

func IsMonotonic(priorities []uint, divider divider.Divider, maxQuantity uint) bool {
	priorities = slices.SortedFunc(slices.Values(priorities), priolist.Compare)

	for combination := range priocombs.Iter(priorities) {
		previous := make(map[uint]uint)

		for quantity := range safe.Iter(0, maxQuantity) {
			actual := make(map[uint]uint)

			divider(combination, quantity, actual)

			distributed, err := distrib.Quantity(actual)
			if err != nil {
				return false
			}

			if distributed != quantity {
				return false
			}

			for priority := range previous {
				if actual[priority] < previous[priority] {
					return false
				}
			}

			for priority := range actual {
				if actual[priority] < previous[priority] {
					return false
				}
			}

			previous = actual
		}
	}

	return true
}
