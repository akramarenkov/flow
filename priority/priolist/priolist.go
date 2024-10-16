// Package used to process the priority list.
package priolist

import (
	"cmp"
	"errors"
	"slices"
)

var (
	ErrEmpty               = errors.New("priority list is empty")
	ErrUnsorted            = errors.New("priority list is unsorted")
	ErrZeroPriorityPresent = errors.New("zero priority is present")
)

// Priority list values comparison function.
func Compare(a, b uint) int {
	return cmp.Compare(b, a)
}

// Validates the priority list.
func IsValid(priorities []uint) error {
	if len(priorities) == 0 {
		return ErrEmpty
	}

	if !slices.IsSortedFunc(priorities, Compare) {
		return ErrUnsorted
	}

	if priorities[len(priorities)-1] == 0 {
		return ErrZeroPriorityPresent
	}

	return nil
}
