// Simplified version of the priority discipline that runs handlers on its own.
package simple

import (
	"errors"

	priocore "github.com/akramarenkov/flow/priority"
	"github.com/akramarenkov/flow/priority/priodefs"
)

var (
	ErrHandleEmpty = errors.New("handle function was not specified")
)

// Callback function called in handlers when an data item is received.
type Handle[P priodefs.Prioritized[Type], Type any] func(prioritized P)

// Options of the created discipline.
type Opts[Type any] struct {
	// Determines in what quantity data items distributed among data handlers
	//
	// For equaling use [divider.Fair] divider, for prioritization use [divider.Rate]
	// divider or custom divider
	Divider priodefs.Divider

	// Callback function called in handlers when an data item is received
	Handle Handle[priodefs.Prioritized[Type], Type]

	// Quantity of data handlers between which data items are distributed
	HandlersQuantity uint

	// Input channels of data items. For terminate the discipline it is necessary and
	// sufficient to close all input channels. Preferably input channels should be
	// buffered for performance reasons. Optimal capacity is equal to the quantity of
	// data handlers
	//
	// Map key is a value of priority. Zero priority is not allowed
	Inputs map[uint]<-chan Type
}

// Adds an input channel with the specified priority to the inputs map.
func (opts *Opts[Type]) AddInput(priority uint, channel <-chan Type) error {
	if priority == 0 {
		return priocore.ErrPriorityZero
	}

	if channel == nil {
		return priocore.ErrInputEmpty
	}

	if opts.Inputs == nil {
		opts.Inputs = make(map[uint]<-chan Type)
	}

	if stored := opts.Inputs[priority]; stored != nil {
		return priocore.ErrInputExists
	}

	opts.Inputs[priority] = channel

	return nil
}

func (opts Opts[Type]) isValid() error {
	if opts.Handle == nil {
		return ErrHandleEmpty
	}

	return nil
}

// Simplified priority discipline.
type Discipline[Type any] struct {
	opts Opts[Type]

	core *priocore.Discipline[Type]
}

// Creates and runs discipline.
func New[Type any](opts Opts[Type]) (*Discipline[Type], error) {
	if err := opts.isValid(); err != nil {
		return nil, err
	}

	core, err := priocore.New(
		priocore.Opts[Type]{
			Divider:          opts.Divider,
			HandlersQuantity: opts.HandlersQuantity,
			Inputs:           opts.Inputs,
		},
	)
	if err != nil {
		return nil, err
	}

	dsc := &Discipline[Type]{
		opts: opts,

		core: core,
	}

	dsc.main()

	return dsc, nil
}

// Returns a channel with errors. If an error occurs (the value from the channel
// is not equal to nil) the discipline terminates its work.
//
// The single nil value means that the discipline has terminated in normal mode:
// after closing and emptying all input channels.
//
// The only place where the error can occurs is the divider. If you are sure that the
// divider is working correctly and the configuration used will not cause an error
// in it, then you are not obliged to read from this channel and you are not obliged
// to check the received value.
func (dsc *Discipline[Type]) Err() <-chan error {
	return dsc.core.Err()
}

func (dsc *Discipline[Type]) main() {
	for range dsc.opts.HandlersQuantity {
		go dsc.handler()
	}
}

func (dsc *Discipline[Type]) handler() {
	for prioritized := range dsc.core.Output() {
		dsc.opts.Handle(prioritized)
		dsc.core.Release(prioritized.Priority)
	}
}
