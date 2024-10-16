// Here are implemented several dividers that determine in what quantity data handlers
// are distributed by priorities.
package divider

import (
	"github.com/akramarenkov/safe"
)

// Determines in what quantity data handlers are distributed by priorities.
//
// Slice of priorities is passed to this function sorted in descending order,
// there can't be zero length and cannot contain zero priority. Divider must not
// change this slice.
//
// Distribution map passed to this function is never nil.
//
// Sum of the quantities of data handlers distributed by this function must be equal
// to the passed quantity of data handlers.
type Divider func(priorities []uint, quantity uint, distribution map[uint]uint)

// Distributes data handlers evenly by priorities.
//
// Used for equaling.
//
// Example results:
//
//   - 6 / [3 2 1] = map[3:2, 2:2, 1:2]
//   - 100 / [70 20 10] = map[70:34, 20:33, 10:33]
func Fair(priorities []uint, quantity uint, distribution map[uint]uint) {
	divider := uint(len(priorities))
	base := quantity / divider
	remainder := quantity % divider

	for _, priority := range priorities {
		part := base

		// Max value of remainder is len(priorities) - 1, so we simply increase
		// distribution by one
		if remainder != 0 {
			part++
			remainder--
		}

		distribution[priority] += part
	}
}

// Distributes data handlers by priorities in proportion to the priority value.
//
// Used for prioritization.
//
// Example results:
//
//   - 6 / [3 2 1] = map[3:3, 2:2, 1:1]
//   - 100 / [70 20 10] = map[70:70, 20:20, 10:10]
func Rate(priorities []uint, quantity uint, distribution map[uint]uint) {
	divider, err := safe.AddMU(priorities...)
	if err != nil {
		panic(err)
	}

	base := quantity / divider
	remainder := quantity % divider

	// Ensures rapid achievement of full distribution filling at small values ​​of
	// quantity to speed up the work of discipline with a small number of handlers
	quicking := uint(len(priorities))

	if quicking > remainder {
		quicking = remainder
	}

	remainder -= quicking

	for _, priority := range priorities {
		part := base * priority

		if quicking != 0 {
			part++
			quicking--
		}

		// Minus one due to quicking
		rating := priority - 1

		if remainder < rating {
			part += remainder
			remainder = 0
		} else {
			part += rating
			remainder -= rating
		}

		distribution[priority] += part
	}
}
