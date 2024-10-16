// Internal package used to create combinations of priorities.
package priocombs

import (
	"iter"
	"math/big"

	"github.com/akramarenkov/safe"
	"github.com/akramarenkov/safe/intspec"
)

// Calculates the quantity of priorities combinations.
//
// If the number of combinations for n priorities is m, then for n+1 priorities the
// number of combinations is m + m + 1. It is easy to see that such an increment
// corresponds to the function 2^n - 1.
func Quantity(priorities []uint) *big.Int {
	quantity := new(big.Int).Sub(
		new(big.Int).Lsh(
			big.NewInt(1),
			uint(len(priorities)),
		),
		big.NewInt(1),
	)

	return quantity
}

// Calculates the size of priorities combinations.
//
// The returned value is intended to be used in the make call as the size parameter so
// is truncated to the maximum value for uint64 if the calculated value exceeds it.
func Size(priorities []uint) uint64 {
	pow, err := safe.Shift[uint64](1, len(priorities))
	if err != nil {
		return intspec.MaxUint64
	}

	return pow - 1
}

// A range iterator by priorities combinations. The returned slice is valid only for
// current iteration of the loop.
//
// It is assumed that priorities are sorted in descending order and are unique, i.e.
// priority values ​​are not taken into account when forming combinations.
func Iter(priorities []uint) iter.Seq[[]uint] {
	iterator := func(yield func([]uint) bool) {
		if len(priorities) == 0 {
			return
		}

		combination := make([]uint, len(priorities))
		shifts := make([]int, len(priorities))
		ids := make([]int, len(priorities))

		level := 0

	stacking: // faster by about 1.5 times than using the function
		for level != -1 {
			remainder := priorities[level+shifts[level]+ids[level]:]

			for _, priority := range remainder {
				combination[level] = priority

				if !yield(combination[:level+1]) {
					return
				}

				if len(remainder) != 1 {
					shifts[level+1] = shifts[level] + ids[level]
					ids[level]++

					level++

					continue stacking
				}
			}

			ids[level] = 0
			level--
		}
	}

	return iterator
}
