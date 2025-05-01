// Discipline used to accumulate items of slices from an input channel into a one slice
// and write that slice to an output channel when the maximum slice size or timeout for
// its accumulation is reached. It works like a join discipline but accepts slices as
// input and unite their items into one slice. Along with this, the input slices are not
// divided between the output slices.
package unite

import (
	"errors"
	"slices"
	"time"
)

var (
	ErrInputEmpty   = errors.New("input channel was not specified")
	ErrJoinSizeZero = errors.New("join size is zero")
)

// Options of the created discipline.
type Opts[Type any] struct {
	// Input channel of data items. For terminate the discipline it is necessary and
	// sufficient to close the input channel. Preferably input channel should be
	// buffered for performance reasons. Optimal capacity is in the range of 1 to 3
	// size of join
	Input <-chan []Type

	// Maximum size of the output slice. Actual size of the output slice may be
	// smaller due to the timeout or closure of the input channel and the fact
	// that the input slices accumulate entirely. Also, the actual size of the output
	// slice may be larger if an slice larger than the maximum size is received at
	// the input
	JoinSize uint

	// By default, to the output channel is written a copy of the accumulated slice.
	// If the NoCopy is set to true, then to the output channel will be directly
	// written the accumulated slice. In this case, after the accumulated slice is
	// no longer used it is necessary to inform the discipline about it by calling
	// Release method
	NoCopy bool

	// Timeout value for output slice accumulation. If the output slice has not been
	// filled completely in the allotted time, then it will be written to the output
	// channel with the data items accumulated during this time. A zero or negative
	// value means that discipline will wait for the missing data items until they
	// appear or the channel is closed (in this case, the accumulated slice will be
	// written to the output channel)
	Timeout time.Duration
}

func (opts Opts[Type]) isValid() error {
	if opts.Input == nil {
		return ErrInputEmpty
	}

	if opts.JoinSize == 0 {
		return ErrJoinSizeZero
	}

	return nil
}

func (opts Opts[Type]) normalize() Opts[Type] {
	if opts.Timeout < 0 {
		opts.Timeout = 0
	}

	return opts
}

// Unite discipline.
type Discipline[Type any] struct {
	opts Opts[Type]

	join    []Type
	output  chan []Type
	release chan struct{}
	timer   *time.Timer
}

// Creates and runs discipline.
func New[Type any](opts Opts[Type]) (*Discipline[Type], error) {
	if err := opts.isValid(); err != nil {
		return nil, err
	}

	opts = opts.normalize()

	dsc := &Discipline[Type]{
		opts: opts,

		join: make([]Type, 0, opts.JoinSize),

		// Value returned by the cap() function is always positive and, in the case of
		// integer overflow due to adding one, the resulting value can only become
		// negative, which will cause a panic when executing make() as same as when
		// specifying a large positive value
		output:  make(chan []Type, 1+cap(opts.Input)),
		release: make(chan struct{}),
	}

	go dsc.main()

	return dsc, nil
}

// Returns output channel.
//
// If this channel is closed, it means that the discipline is terminated.
func (dsc *Discipline[Type]) Output() <-chan []Type {
	return dsc.output
}

// Marks output slice as no longer used outside of the discipline.
//
// Must be called, if NoCopy option is set to true, after the output slice is
// no longer used outside of the discipline. However, calling this method is also
// possible if the NoCopy option is set to false.
func (dsc *Discipline[Type]) Release() {
	if dsc.opts.NoCopy {
		dsc.release <- struct{}{}
	}
}

func (dsc *Discipline[Type]) main() {
	defer close(dsc.output)
	defer close(dsc.release)

	if dsc.opts.Timeout == 0 {
		dsc.loopWithoutTimeout()
		return
	}

	dsc.loop()
}

func (dsc *Discipline[Type]) loopWithoutTimeout() {
	defer dsc.pass()

	for item := range dsc.opts.Input {
		dsc.add(item)
	}
}

func (dsc *Discipline[Type]) loop() {
	dsc.timer = time.NewTimer(dsc.opts.Timeout)
	defer dsc.timer.Stop()

	defer dsc.pass()

	for {
		select {
		case <-dsc.timer.C:
			dsc.pass()
		case item, opened := <-dsc.opts.Input:
			if !opened {
				return
			}

			dsc.add(item)
		}
	}
}

func (dsc *Discipline[Type]) add(item []Type) {
	if uint(len(item)) >= dsc.opts.JoinSize {
		dsc.pass()
		dsc.forward(item)

		return
	}

	// Integer overflow is impossible because len() function returns only positive
	// values for the int type and the sum of the two maximum values for the int type is
	// less than the maximum value for the uint type by one
	if uint(len(item))+uint(len(dsc.join)) > dsc.opts.JoinSize {
		dsc.pass()
	}

	dsc.join = append(dsc.join, item...)

	// Integer overflow is impossible because len() function returns only positive
	// values for type int and the maximum value for type int is less than the
	// maximum value for type uint
	if uint(len(dsc.join)) < dsc.opts.JoinSize {
		return
	}

	dsc.pass()
}

func (dsc *Discipline[Type]) pass() {
	if len(dsc.join) == 0 {
		// defer statement is not used to allow inlining of the current function
		dsc.resetTimer()
		return
	}

	dsc.send(dsc.join)
	dsc.resetJoin()
	dsc.resetTimer()
}

func (dsc *Discipline[Type]) forward(item []Type) {
	dsc.send(item)
	dsc.resetTimer()
}

func (dsc *Discipline[Type]) send(item []Type) {
	item = dsc.prepareItem(item)

	dsc.output <- item

	if dsc.opts.NoCopy {
		<-dsc.release
	}
}

func (dsc *Discipline[Type]) prepareItem(item []Type) []Type {
	if dsc.opts.NoCopy {
		return item
	}

	return slices.Clone(item)
}

func (dsc *Discipline[Type]) resetJoin() {
	dsc.join = dsc.join[:0]
}

func (dsc *Discipline[Type]) resetTimer() {
	if dsc.opts.Timeout == 0 {
		return
	}

	dsc.timer.Reset(dsc.opts.Timeout)
}
