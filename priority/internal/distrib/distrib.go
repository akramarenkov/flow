// Internal package used to calculate distribution parameters.
package distrib

import "github.com/akramarenkov/safe"

// Calculates the total quantity in the distribution.
func Quantity(distribution map[uint]uint) (uint, error) {
	quantity := uint(0)

	for _, amount := range distribution {
		sum, err := safe.AddU(quantity, amount)
		if err != nil {
			return 0, err
		}

		quantity = sum
	}

	return quantity, nil
}

// Checks that there are no zero values ​​in the distribution and the distribution
// itself is not empty.
func IsFilled(distribution map[uint]uint) bool {
	for _, quantity := range distribution {
		if quantity == 0 {
			return false
		}
	}

	return len(distribution) != 0
}
