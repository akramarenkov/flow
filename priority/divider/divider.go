// Here are implemented several dividers that determine in what quantity data items
// distributed among data handlers.
package divider

import (
	"fmt"

	"github.com/akramarenkov/safe"
)

// Distributes data items evenly among data handlers.
//
// Used for equaling.
//
// To keep processing data items of all priorities, the quantity of data handlers must
// be no less than the quantity of priorities.
//
// Example results:
//
//   - 6 / [3 2 1] = map[3:2 2:2 1:2]
//   - 10 / [7 2 1] = map[7:4 2:3 1:3]
//   - 100 / [70 20 10] = map[70:34 20:33 10:33]
func Fair(quantity uint, priorities []uint, distribution map[uint]uint) error {
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

	return nil
}

// Distributes data items among data handlers in ratio to the priority of data item.
//
// Used for prioritization.
//
// To keep processing data items of all priorities, the quantity of data handlers must
// be no less than the sum of the values ​​of all priorities.
//
// Example results:
//
//   - 6 / [3 2 1] = map[3:3 2:2 1:1]
//   - 10 / [7 2 1] = map[7:7 2:2 1:1]
//   - 100 / [70 20 10] = map[70:70 20:20 10:10]
func Rate(quantity uint, priorities []uint, distribution map[uint]uint) error {
	divider, err := safe.AddMU(priorities...)
	if err != nil {
		return fmt.Errorf("calculation of the sum of priorities: %w", err)
	}

	base := quantity / divider
	remainder := quantity % divider

	// Provides rapid achievement of full distribution filling at small values ​​of
	// quantity of data handlers to speed up work of the discipline
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

	return nil
}
