// Internal package used to calculate distribution parameters.
package distrib

import "github.com/akramarenkov/safe"

// Calculates the total quantity in the distribution for specified priorities.
func Quantity(priorities []uint, distribution map[uint]uint) (uint, error) {
	if len(distribution) == 0 {
		return 0, nil
	}

	quantity := uint(0)

	for _, priority := range priorities {
		sum, err := safe.AddU(quantity, distribution[priority])
		if err != nil {
			return 0, err
		}

		quantity = sum
	}

	return quantity, nil
}

// Checks that there are no zero values ​​in the distribution for specified priorities and
// that the distribution and priority list themselves are not empty.
func IsFilled(priorities []uint, distribution map[uint]uint) bool {
	if len(distribution) == 0 {
		return false
	}

	for _, priority := range priorities {
		if distribution[priority] == 0 {
			return false
		}
	}

	return len(priorities) != 0
}
