package measuring

import "errors"

var (
	ErrDataMixedBetweenPriorities = errors.New("data is mixed between priorities")
	ErrDataPartiallyLost          = errors.New("data is partially lost")
	ErrHandlersQuantityZero       = errors.New("quantity of data handlers is zero")
	ErrMeasureDuplicated          = errors.New("measure is duplicated")
	ErrMeasurementsIsIncomplete   = errors.New("measurements is incomplete")
	ErrUnexpectedMeasureKind      = errors.New("unexpected kind of measure is received")
)
