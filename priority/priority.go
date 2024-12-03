// Discipline used to distribute data items between handlers in quantity
// corresponding to the priority of the data items.
package priority

import (
	"slices"
	"time"

	"github.com/akramarenkov/flow/priority/internal/distrib"
	"github.com/akramarenkov/flow/priority/types"
)

const (
	defaultIdleDelay = 1 * time.Nanosecond
)

// Options of the created discipline.
type Opts[Type any] struct {
	// Determines in what quantity data items distributed among data handlers
	//
	// For equaling use divider.Fair divider, for prioritization use divider.Rate
	// divider or custom divider
	Divider types.Divider
	// Quantity of data handlers between which data items are distributed
	HandlersQuantity uint
	// Input channels of data items. For terminate the discipline it is necessary and
	// sufficient to close all input channels. Preferably input channels should be
	// buffered for performance reasons. Optimal capacity is in the range of 1 to 3
	// times of quantity of data handlers
	//
	// Map key is a value of priority
	Inputs map[uint]<-chan Type
}

// Adds an input channel with the specified priority to the inputs map.
func (opts *Opts[Type]) AddInput(priority uint, channel <-chan Type) error {
	if priority == 0 {
		return ErrPriorityZero
	}

	if channel == nil {
		return ErrInputEmpty
	}

	if opts.Inputs == nil {
		opts.Inputs = make(map[uint]<-chan Type)
	}

	if stored := opts.Inputs[priority]; stored != nil {
		return ErrInputExists
	}

	opts.Inputs[priority] = channel

	return nil
}

func (opts Opts[Type]) isInputEmpty() error {
	if len(opts.Inputs) == 0 {
		return ErrInputEmpty
	}

	for _, channel := range opts.Inputs {
		if channel == nil {
			return ErrInputEmpty
		}
	}

	return nil
}

func (opts Opts[Type]) isValid() error {
	if opts.Divider == nil {
		return ErrDividerEmpty
	}

	if opts.HandlersQuantity == 0 {
		return ErrHandlersQuantityZero
	}

	if err := opts.isInputEmpty(); err != nil {
		return err
	}

	if channel := opts.Inputs[0]; channel != nil {
		return ErrPriorityZero
	}

	return nil
}

func (opts Opts[Type]) disciplineInputs() map[uint]input[Type] {
	inputs := make(map[uint]input[Type], len(opts.Inputs))

	for priority, channel := range opts.Inputs {
		input := input[Type]{
			Channel: channel,
		}

		inputs[priority] = input
	}

	return inputs
}

// Priority discipline.
type Discipline[Type any] struct {
	opts Opts[Type]

	inputs  map[uint]input[Type]
	output  chan types.Prioritized[Type]
	release chan uint

	// priority list corresponding to all input channels - main priority list
	priorities []uint
	// priority list whose actual distribution did not reach strategic
	unachieved []uint
	// priority list whose actual distribution did not reach operative
	unreached []uint
	// priority list from whose channels it managed to get all data items at
	// ​​input/output stage for priorities from the unachieved list and, since the
	// unachieved list may not be complete with respect to the main priority list
	// then, at any previous input/output stages - interim main priority list
	useful []uint

	// actual distribution of data items by priorities
	actual map[uint]uint
	// distribution of data items filled by useful priority list and total quantity of
	// data handlers - interim strategic distribution
	operative map[uint]uint
	// distribution of data items filled by main priority list and total quantity of
	// data handlers
	strategic map[uint]uint
	// distribution on whose quantities ​​input/output is performed
	tactic map[uint]uint

	err chan error
}

// Creates and runs discipline.
func New[Type any](opts Opts[Type]) (*Discipline[Type], error) {
	if err := opts.isValid(); err != nil {
		return nil, err
	}

	inputs, priorities, strategic, err := prepare(opts)
	if err != nil {
		return nil, err
	}

	dsc := &Discipline[Type]{
		opts: opts,

		inputs:  inputs,
		output:  make(chan types.Prioritized[Type], opts.HandlersQuantity),
		release: make(chan uint, opts.HandlersQuantity),

		priorities: priorities,

		actual:    make(map[uint]uint),
		operative: make(map[uint]uint),
		strategic: strategic,
		tactic:    make(map[uint]uint),

		err: make(chan error, 1),
	}

	go dsc.main()

	return dsc, nil
}

func prepare[Type any](opts Opts[Type]) (
	map[uint]input[Type],
	[]uint,
	map[uint]uint,
	error,
) {
	inputs := opts.disciplineInputs()

	priorities := make([]uint, 0, len(inputs))
	strategic := make(map[uint]uint, len(inputs))

	for priority := range inputs {
		priorities = append(priorities, priority)
	}

	slices.SortFunc(priorities, Compare)

	err := divide(opts.Divider, opts.HandlersQuantity, priorities, strategic)
	if err != nil {
		return nil, nil, nil, err
	}

	if !distrib.IsFilled(priorities, strategic) {
		return nil, nil, nil, ErrHandlersQuantityTooSmall
	}

	return inputs, priorities, strategic, nil
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
// Handlers must call this method after the current data item has been processed.
func (dsc *Discipline[Type]) Release(priority uint) {
	dsc.release <- priority
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
	return dsc.err
}

func (dsc *Discipline[Type]) main() {
	defer close(dsc.err)
	defer close(dsc.output)
	defer close(dsc.release)

	if err := dsc.loop(); err != nil {
		dsc.err <- err
	}
}

func (dsc *Discipline[Type]) loop() error {
	defer dsc.waitFullReleased()

	for {
		distributed, err := dsc.distribute()
		if err != nil {
			return err
		}

		dsc.collectReleases()

		if distributed == 0 {
			if dsc.isInputsClosed() {
				return nil
			}

			time.Sleep(defaultIdleDelay)
		}
	}
}

func (dsc *Discipline[Type]) waitFullReleased() {
	for !dsc.isFullyReleased() {
		dsc.waitRelease()
	}
}

func (dsc *Discipline[Type]) isFullyReleased() bool {
	for _, priority := range dsc.priorities {
		if dsc.actual[priority] != 0 {
			return false
		}
	}

	return true
}

func (dsc *Discipline[Type]) collectReleases() {
	for range len(dsc.release) {
		dsc.waitRelease()
	}
}

func (dsc *Discipline[Type]) waitRelease() {
	dsc.actual[<-dsc.release]--
}

func (dsc *Discipline[Type]) isInputsClosed() bool {
	for _, input := range dsc.inputs {
		if !input.Closed {
			return false
		}
	}

	return true
}

func (dsc *Discipline[Type]) distribute() (uint, error) {
	distributed := uint(0)

	if err := dsc.waitFillingUnachieved(); err != nil {
		return distributed, err
	}

	distributed += dsc.transfer(dsc.unachieved)

	proceed, err := dsc.fillOperative()
	if err != nil {
		return distributed, err
	}

	if !proceed {
		return distributed, nil
	}

	if err := dsc.waitFillingUnreached(); err != nil {
		return distributed, err
	}

	distributed += dsc.transfer(dsc.unreached)

	return distributed, nil
}

func (dsc *Discipline[Type]) waitFillingUnachieved() error {
	for {
		filled, err := dsc.fillUnachieved()
		if err != nil {
			return err
		}

		if filled {
			return nil
		}

		dsc.waitRelease()
	}
}

// Fills tactical distribution for unachieved priorities.
func (dsc *Discipline[Type]) fillUnachieved() (bool, error) {
	vacant := dsc.vacantHandlers()

	if vacant == 0 {
		return false, nil
	}

	dsc.prepareUnachieved()
	dsc.resetTactic()

	if err := divide(dsc.opts.Divider, vacant, dsc.unachieved, dsc.tactic); err != nil {
		return false, err
	}

	return distrib.IsFilled(dsc.unachieved, dsc.tactic), nil
}

func (dsc *Discipline[Type]) vacantHandlers() uint {
	// integer overflow or incorrect counting are not possible here because
	// the correctness of the distribution is checked at each dividing
	return dsc.opts.HandlersQuantity - dsc.busyHandlers()
}

func (dsc *Discipline[Type]) busyHandlers() uint {
	busy := uint(0)

	// integer overflow or incorrect counting are not possible here because
	// the correctness of the distribution is checked at each dividing

	for _, priority := range dsc.priorities {
		busy += dsc.actual[priority]
	}

	return busy
}

func (dsc *Discipline[Type]) prepareUnachieved() {
	dsc.unachieved = dsc.unachieved[:0]

	for _, priority := range dsc.priorities {
		if dsc.actual[priority] < dsc.strategic[priority] {
			dsc.unachieved = append(dsc.unachieved, priority)
		}
	}
}

func (dsc *Discipline[Type]) resetTactic() {
	for _, priority := range dsc.priorities {
		dsc.tactic[priority] = 0
	}
}

func (dsc *Discipline[Type]) fillOperative() (bool, error) {
	dsc.prepareUseful()

	if len(dsc.useful) == 0 {
		return false, nil
	}

	dsc.resetOperative()

	err := divide(dsc.opts.Divider, dsc.opts.HandlersQuantity, dsc.useful, dsc.operative)
	if err != nil {
		return false, err
	}

	if !distrib.IsFilled(dsc.useful, dsc.operative) {
		return false, ErrDividerBad
	}

	return true, nil
}

func (dsc *Discipline[Type]) prepareUseful() {
	dsc.useful = dsc.useful[:0]

	// is used the main priority list because it is necessary to take into account
	// the results of not only the current stage of input/output, but also the previous
	// ones
	for _, priority := range dsc.priorities {
		if dsc.tactic[priority] == 0 {
			dsc.useful = append(dsc.useful, priority)
		}
	}
}

func (dsc *Discipline[Type]) resetOperative() {
	for _, priority := range dsc.priorities {
		dsc.operative[priority] = 0
	}
}

func (dsc *Discipline[Type]) waitFillingUnreached() error {
	for {
		filled, err := dsc.fillUnreached()
		if err != nil {
			return err
		}

		if filled {
			return nil
		}

		dsc.waitRelease()
	}
}

func (dsc *Discipline[Type]) fillUnreached() (bool, error) {
	vacant := dsc.vacantHandlers()

	if vacant == 0 {
		return false, nil
	}

	dsc.prepareUnreached()
	dsc.resetTactic()

	if err := divide(dsc.opts.Divider, vacant, dsc.unreached, dsc.tactic); err != nil {
		return false, err
	}

	return distrib.IsFilled(dsc.unreached, dsc.tactic), nil
}

func (dsc *Discipline[Type]) prepareUnreached() {
	dsc.unreached = dsc.unreached[:0]

	for _, priority := range dsc.useful {
		if dsc.actual[priority] < dsc.operative[priority] {
			dsc.unreached = append(dsc.unreached, priority)
		}
	}
}

func (dsc *Discipline[Type]) transfer(priorities []uint) uint {
	transferred := uint(0)

	for _, priority := range priorities {
		if dsc.inputs[priority].Closed {
			continue
		}

		transferred += dsc.pass(priority)
	}

	return transferred
}

func (dsc *Discipline[Type]) pass(priority uint) uint {
	passed := uint(0)

	for dsc.tactic[priority] != 0 {
		select {
		case item, opened := <-dsc.inputs[priority].Channel:
			if !opened {
				dsc.markInputAsClosed(priority)
				return passed
			}

			passed += dsc.send(item, priority)
		default:
			return passed
		}
	}

	return passed
}

func (dsc *Discipline[Type]) markInputAsClosed(priority uint) {
	input := dsc.inputs[priority]

	input.Closed = true

	dsc.inputs[priority] = input
}

func (dsc *Discipline[Type]) send(item Type, priority uint) uint {
	prioritized := types.Prioritized[Type]{
		Item:     item,
		Priority: priority,
	}

	dsc.output <- prioritized

	dsc.tactic[priority]--
	dsc.actual[priority]++

	return 1
}
