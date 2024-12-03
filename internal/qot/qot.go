// Internal package describing the quantity over time structure.
package qot

import (
	"time"
)

// Quantity over time.
type QoT struct {
	Quantity uint
	Time     time.Duration
}
