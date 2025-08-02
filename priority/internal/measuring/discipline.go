package measuring

import (
	"github.com/akramarenkov/flow/priority/priodefs"
)

// Interface of priority discipline.
type Discipline[Type any] interface {
	Err() <-chan error
	Output() <-chan priodefs.Prioritized[Type]
	Release(priority uint)
}
