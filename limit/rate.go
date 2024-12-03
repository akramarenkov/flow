package limit

import (
	"errors"
	"time"
)

var (
	ErrIntervalNegative = errors.New("interval is negative")
	ErrIntervalZero     = errors.New("interval is zero")
	ErrQuantityZero     = errors.New("quantity is zero")
)

// Quantity of data items passed per time Interval.
type Rate struct {
	Interval time.Duration
	Quantity uint64
}

// Validates field values. Interval cannot be negative or equal to zero. Quantity
// cannot be equal to zero.
func (rt Rate) IsValid() error {
	if rt.Interval < 0 {
		return ErrIntervalNegative
	}

	if rt.Interval == 0 {
		return ErrIntervalZero
	}

	if rt.Quantity == 0 {
		return ErrQuantityZero
	}

	return nil
}
