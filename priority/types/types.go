// Data types of the priority discipline.
package types

// Determines in what quantity data items of the specified priorities are
// distributed among the specified quantity of data handlers.
//
// Priority list passed to the divider is always sorted in descending order,
// cannot be of zero length and cannot contain a zero priority. The divider must not
// modify this slice.
//
// Distribution map passed to the divider cannot be nil.
//
// Total quantity of data items in the distribution created by the divider
// must be equal to the quantity of data handlers passed to the divider.
//
// Quantity of data items for each priority in the distribution
// created by the divider must be monotonically non-decreasing as the quantity of data
// handlers passed to the divider increases.
type Divider func(quantity uint, priorities []uint, distribution map[uint]uint) error

// Describes the data item distributed by the priority discipline.
type Prioritized[Type any] struct {
	Item     Type
	Priority uint
}
