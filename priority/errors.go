package priority

import "errors"

var (
	ErrDividerBad               = errors.New("divider creates an incorrect distribution")
	ErrDividerEmpty             = errors.New("divider was not specified")
	ErrHandlersQuantityTooSmall = errors.New("quantity of data handlers is too small")
	ErrHandlersQuantityZero     = errors.New("quantity of data handlers is zero")
	ErrInputEmpty               = errors.New("input channel was not specified")
	ErrInputExists              = errors.New("input channel already specified")
	ErrPriorityZero             = errors.New("zero priority is specified")
)
