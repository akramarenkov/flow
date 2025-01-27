package measuring

import (
	"sync"
	"time"

	"github.com/akramarenkov/breaker/closing"
	"github.com/akramarenkov/flow/priority/types"
	"github.com/akramarenkov/span"
	"github.com/akramarenkov/starter"
)

const (
	waitDevastationDelay = 1 * time.Nanosecond
)

// Type of action performed by the measurer.
type actionKind int

const (
	actionKindDelay actionKind = iota + 1
	actionKindWaitDevastation
	actionKindWrite
	actionKindWriteWithDelay
)

// Action performed by the measurer.
type action struct {
	Delay    time.Duration
	Kind     actionKind
	Quantity uint
}

// Used to measuring discipline.
type Measurer struct {
	handlersQuantity uint
	inputCapacities  []uint

	actions   map[uint][]action
	channels  map[uint]chan uint
	durations map[uint]time.Duration

	breaker      *closing.Closing
	measurements []Measure
	measuring    chan Measure
	spans        map[uint]span.Span[uint]
	starter      *starter.Starter
	wg           *sync.WaitGroup
}

// Creates Measurer instance.
//
// If the capacity of the input (for the discipline) channels is not specified,
// then the capacity of each channel will be equal to the quantity of data handlers.
//
// When specifying multiple input channel capacities, the value of the first
// one will be used.
func NewMeasurer(handlersQuantity uint, inputCapacity ...uint) (*Measurer, error) {
	if handlersQuantity == 0 {
		return nil, ErrHandlersQuantityZero
	}

	msr := &Measurer{
		handlersQuantity: handlersQuantity,
		inputCapacities:  append([]uint(nil), inputCapacity...),

		actions:   make(map[uint][]action),
		channels:  make(map[uint]chan uint),
		durations: make(map[uint]time.Duration),

		wg: new(sync.WaitGroup),
	}

	return msr, nil
}

// Returns the quantity of data handlers specified when the instance was created.
func (msr *Measurer) HandlersQuantity() uint {
	return msr.handlersQuantity
}

// Adds to the actions list a write of the specified quantity of data items to the input
// (for the discipline) channel of the specified priority.
func (msr *Measurer) AddWrite(priority uint, quantity uint) {
	act := action{
		Kind:     actionKindWrite,
		Quantity: quantity,
	}

	msr.addAction(priority, act)
}

// Adds to the actions list a write of the specified quantity of data items to the input
// (for the discipline) channel of the specified priority.
//
// Before writing each data item, the specified delay occurs.
func (msr *Measurer) AddWriteWithDelay(priority uint, quantity uint, delay time.Duration) {
	act := action{
		Delay:    delay,
		Kind:     actionKindWriteWithDelay,
		Quantity: quantity,
	}

	msr.addAction(priority, act)
}

// Adds to the actions list a waiting for the input (for the discipline) channel to
// be devastated.
func (msr *Measurer) AddWaitDevastation(priority uint) {
	act := action{
		Kind: actionKindWaitDevastation,
	}

	msr.addAction(priority, act)
}

// Adds to the actions list a delay in further execution of actions.
func (msr *Measurer) AddDelay(priority uint, delay time.Duration) {
	act := action{
		Kind:  actionKindDelay,
		Delay: delay,
	}

	msr.addAction(priority, act)
}

func (msr *Measurer) addAction(priority uint, action action) {
	msr.actions[priority] = append(msr.actions[priority], action)
}

// Returns the total quantity of data items that will be written to the input
// (for the discipline) channels.
func (msr *Measurer) itemsQuantity() uint {
	quantity := uint(0)

	for _, actions := range msr.actions {
		for _, action := range actions {
			switch action.Kind {
			case actionKindWrite, actionKindWriteWithDelay:
				quantity += action.Quantity
			}
		}
	}

	return quantity
}

// Returns the quantity of data items that will be written to the input (for the
// discipline) channel of specified priority.
func (msr *Measurer) itemsQuantityForPriority(priority uint) uint {
	quantity := uint(0)

	for _, action := range msr.actions[priority] {
		switch action.Kind {
		case actionKindWrite, actionKindWriteWithDelay:
			quantity += action.Quantity
		}
	}

	return quantity
}

// Sets the processing duration of one data item received from the discipline.
func (msr *Measurer) SetProcessingDuration(priority uint, duration time.Duration) {
	msr.durations[priority] = duration
}

// Recreates the input (for the discipline) channels where data items will be written.
//
// Must be called before [Measurer.Play].
func (msr *Measurer) Inputs() map[uint]<-chan uint {
	clear(msr.channels)

	inputs := make(map[uint]<-chan uint, len(msr.actions))

	for priority := range msr.actions {
		channel := make(chan uint, msr.inputCapacity())

		msr.channels[priority] = channel
		inputs[priority] = channel
	}

	return inputs
}

func (msr *Measurer) inputCapacity() uint {
	if len(msr.inputCapacities) != 0 {
		return msr.inputCapacities[0]
	}

	return msr.handlersQuantity
}

// Writes data items to the input (for the discipline) channels, reads data items from
// the output channel of the discipline and collects measurements about the
// time it takes to receive, process and release data items.
//
// Must be called after [Measurer.Inputs].
func (msr *Measurer) Play(discipline Discipline[uint]) ([]Measure, error) {
	msr.prepare()
	msr.handlers(discipline)
	msr.writers()

	if err := <-discipline.Err(); err != nil {
		msr.breaker.Close()
		msr.wg.Wait()
		close(msr.measuring)

		return nil, err
	}

	msr.wg.Wait()
	close(msr.measuring)
	msr.breaker.Close()

	for measure := range msr.measuring {
		msr.measurements = append(msr.measurements, measure)
	}

	if err := isCorrectMeasurements(msr.measurements, msr.spans); err != nil {
		return nil, err
	}

	return msr.measurements, nil
}

func (msr *Measurer) prepare() {
	quantity := msr.measurementsQuantity()

	msr.breaker = closing.New()
	msr.measurements = make([]Measure, 0, quantity)
	msr.measuring = make(chan Measure, quantity)
	msr.starter = starter.New()

	msr.prepareSpans()
}

// Returns the quantity of measure items that should will be collected.
func (msr *Measurer) measurementsQuantity() uint {
	return measurementsPerDataItem * msr.itemsQuantity()
}

// Used to get unique values ​​of data items for all priorities.
func (msr *Measurer) prepareSpans() {
	begin := uint(0)

	msr.spans = make(map[uint]span.Span[uint], len(msr.channels))

	for priority := range msr.channels {
		quantity := msr.itemsQuantityForPriority(priority)

		if quantity == 0 {
			continue
		}

		msr.spans[priority] = span.Span[uint]{
			Begin: begin,
			End:   begin + quantity - 1,
		}

		begin += quantity
	}
}

func (msr *Measurer) writers() {
	for priority := range msr.channels {
		msr.wg.Add(1)

		go msr.writer(priority, msr.spans[priority].Begin)
	}
}

func (msr *Measurer) writer(priority uint, shift uint) {
	defer msr.wg.Done()
	defer close(msr.channels[priority])

	sequence := shift

	for _, action := range msr.actions[priority] {
		switch action.Kind {
		case actionKindWrite, actionKindWriteWithDelay:
			increased, stop := msr.write(action, msr.channels[priority], sequence)
			if stop {
				return
			}

			sequence = increased
		case actionKindWaitDevastation:
			if stop := msr.waitDevastation(msr.channels[priority]); stop {
				return
			}
		case actionKindDelay:
			time.Sleep(action.Delay)
		}
	}
}

func (msr *Measurer) write(action action, channel chan uint, sequence uint) (uint, bool) {
	for range action.Quantity {
		select {
		case <-msr.breaker.IsClosed():
			return sequence, true
		case channel <- sequence:
		}

		if action.Kind == actionKindWriteWithDelay {
			time.Sleep(action.Delay)
		}

		sequence++
	}

	return sequence, false
}

func (msr *Measurer) waitDevastation(channel chan uint) bool {
	ticker := time.NewTicker(waitDevastationDelay)
	defer ticker.Stop()

	for {
		select {
		case <-msr.breaker.IsClosed():
			return true
		case <-ticker.C:
			if len(channel) == 0 {
				return false
			}
		}
	}
}

func (msr *Measurer) handlers(discipline Discipline[uint]) {
	defer msr.starter.Go()

	for range msr.handlersQuantity {
		msr.wg.Add(1)
		msr.starter.Ready()

		go msr.handler(discipline)
	}
}

func (msr *Measurer) handler(discipline Discipline[uint]) {
	defer msr.wg.Done()

	msr.starter.Set()

	for item := range discipline.Output() {
		msr.handle(item, discipline)
	}
}

func (msr *Measurer) handle(item types.Prioritized[uint], discipline Discipline[uint]) {
	received := Measure{
		Item:     item.Item,
		Kind:     KindReceived,
		Priority: item.Priority,
		Time:     time.Since(msr.starter.StartedAt()),
	}

	msr.measuring <- received

	time.Sleep(msr.durations[item.Priority])

	processed := Measure{
		Item:     item.Item,
		Kind:     KindProcessed,
		Priority: item.Priority,
		Time:     time.Since(msr.starter.StartedAt()),
	}

	msr.measuring <- processed

	discipline.Release(item.Priority)

	completed := Measure{
		Item:     item.Item,
		Kind:     KindCompleted,
		Priority: item.Priority,
		Time:     time.Since(msr.starter.StartedAt()),
	}

	msr.measuring <- completed
}
