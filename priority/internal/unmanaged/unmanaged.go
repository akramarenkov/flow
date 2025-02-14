// Internal package with the implementation of a discipline that does not distribute
// data items between handlers in quantity corresponding to the priority of the data
// items.
package unmanaged

import (
	"errors"
	"sync"

	"github.com/akramarenkov/flow/priority/types"

	"github.com/akramarenkov/breaker/closing"
)

var (
	ErrDeliberateFailure    = errors.New("deliberate failure")
	ErrHandlersQuantityZero = errors.New("quantity of data handlers is zero")
	ErrInputEmpty           = errors.New("input channel was not specified")
)

// Options of the created discipline.
type Opts[Type any] struct {
	// Specifies after processing which (by quantity) data item for each priority the
	// discipline should terminate with an error. An empty map means termination in
	// normal mode: without an error, after closing and emptying all input channels
	FailAfter map[uint]uint

	// Quantity of data handlers between which data items are distributed
	HandlersQuantity uint

	// Input channels of data items. For terminate the discipline it is necessary and
	// sufficient to close all input channels. Preferably input channels should be
	// buffered for performance reasons.
	//
	// Map key is a value of priority. Zero priority is not allowed
	Inputs map[uint]<-chan Type

	// Quantity of data items for each priority that will not be written to the output
	// channel
	Misses map[uint]uint
}

func (opts Opts[Type]) isValid() error {
	if opts.HandlersQuantity == 0 {
		return ErrHandlersQuantityZero
	}

	return opts.isEmptyInput()
}

func (opts Opts[Type]) isEmptyInput() error {
	for _, channel := range opts.Inputs {
		if channel != nil {
			return nil
		}
	}

	return ErrInputEmpty
}

// Unmanaged discipline.
type Discipline[Type any] struct {
	opts Opts[Type]

	breaker  *closing.Closing
	err      chan error
	failures chan error
	output   chan types.Prioritized[Type]
	wg       sync.WaitGroup
}

// Creates and runs discipline.
func New[Type any](opts Opts[Type]) (*Discipline[Type], error) {
	if err := opts.isValid(); err != nil {
		return nil, err
	}

	dsc := &Discipline[Type]{
		opts: opts,

		breaker:  closing.New(),
		err:      make(chan error, 1),
		failures: make(chan error, len(opts.Inputs)),
		output:   make(chan types.Prioritized[Type], opts.HandlersQuantity),
	}

	go dsc.main()

	return dsc, nil
}

// Returns output channel.
//
// If this channel is closed, it means that the discipline is terminated.
func (dsc *Discipline[Type]) Output() <-chan types.Prioritized[Type] {
	return dsc.output
}

// Marks that current data item has been processed and handler is ready to receive new
// data item.
//
// Does nothing, defined for compatibility with the priority discipline interface.
func (*Discipline[Type]) Release(uint) {
}

// Returns a channel with errors. If an error occurs (the value from the channel
// is not equal to nil) the discipline terminates its work.
//
// The single nil value means that the discipline has terminated in normal mode:
// after closing and emptying all input channels.
func (dsc *Discipline[Type]) Err() <-chan error {
	return dsc.err
}

func (dsc *Discipline[Type]) main() {
	defer close(dsc.err)
	defer close(dsc.output)

	if err := dsc.loop(); err != nil {
		dsc.err <- err
	}
}

func (dsc *Discipline[Type]) loop() error {
	for priority := range dsc.opts.Inputs {
		dsc.wg.Add(1)

		go dsc.io(priority)
	}

	dsc.wg.Wait()

	close(dsc.failures)

	return <-dsc.failures
}

func (dsc *Discipline[Type]) io(priority uint) {
	defer dsc.wg.Done()

	if dsc.opts.Inputs[priority] == nil {
		return
	}

	if dsc.opts.FailAfter != nil && dsc.opts.FailAfter[priority] != 0 {
		dsc.faulty(priority)
		return
	}

	dsc.unfaulty(priority)
}

func (dsc *Discipline[Type]) faulty(priority uint) {
	processed := uint(0)

	misses := dsc.opts.Misses[priority]

	for {
		select {
		case <-dsc.breaker.IsClosed():
			return
		case item, opened := <-dsc.opts.Inputs[priority]:
			if !opened {
				return
			}

			if misses != 0 {
				misses--
				continue
			}

			dsc.send(item, priority)

			processed++

			if processed == dsc.opts.FailAfter[priority] {
				dsc.fail(ErrDeliberateFailure)
				dsc.breaker.Close()

				return
			}
		}
	}
}

func (dsc *Discipline[Type]) unfaulty(priority uint) {
	misses := dsc.opts.Misses[priority]

	for {
		select {
		case <-dsc.breaker.IsClosed():
			return
		case item, opened := <-dsc.opts.Inputs[priority]:
			if !opened {
				return
			}

			if misses != 0 {
				misses--
				continue
			}

			dsc.send(item, priority)
		}
	}
}

func (dsc *Discipline[Type]) send(item Type, priority uint) {
	prioritized := types.Prioritized[Type]{
		Item:     item,
		Priority: priority,
	}

	dsc.output <- prioritized
}

func (dsc *Discipline[Type]) fail(err error) {
	dsc.failures <- err
}
