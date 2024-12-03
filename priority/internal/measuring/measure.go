package measuring

import (
	"cmp"
	"time"
)

const (
	measurementsPerDataItem = 3
)

// Type of one measuring. Specifies the position of the measurement in the
// processing sequence of the data item received from the discipline.
type Kind int

const (
	KindCompleted Kind = iota + 1
	KindProcessed
	KindReceived
)

// Describes one measuring.
type Measure struct {
	// Data item received from the discipline
	Item uint
	// Position of the measurement in the processing sequence of the data item
	Kind Kind
	// Priority of data item
	Priority uint
	// Relative time of measurement execution
	Time time.Duration
}

// Compare function for sorting measurements by data item.
func CompareItem(first, second Measure) int {
	return cmp.Compare(first.Item, second.Item)
}

// Compare function for sorting measurements by time.
func CompareTime(first, second Measure) int {
	return cmp.Compare(first.Time, second.Time)
}

// Delete function that keeps only measurements of the Received kind.
func KeepReceived(item Measure) bool {
	return item.Kind != KindReceived
}
