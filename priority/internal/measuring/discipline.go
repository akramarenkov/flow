package measuring

import (
	"github.com/akramarenkov/flow/priority/types"
)

// Interface of priority discipline.
type Discipline[Type any] interface {
	Err() <-chan error
	Output() <-chan types.Prioritized[Type]
	Release(priority uint)
}
