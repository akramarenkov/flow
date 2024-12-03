package priority

import "cmp"

// Compare function for sorting priorities.
func Compare(first, second uint) int {
	return cmp.Compare(second, first)
}
