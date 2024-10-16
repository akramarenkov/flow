package measurer

import (
	"github.com/akramarenkov/flow/priority/types"
)

type Discipline[Type any] interface {
	Output() <-chan types.Prioritized[Type]
	Release(priority uint)
	Err() <-chan error
}
