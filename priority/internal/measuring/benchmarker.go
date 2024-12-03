package measuring

import "sync"

// Used to benchmark discipline.
type Benchmarker struct {
	handlersQuantity uint
	inputCapacities  []uint

	channels      map[uint]chan uint
	itemsQuantity map[uint]uint
	wg            *sync.WaitGroup
}

// Creates Benchmarker instance.
//
// If the capacity of the input (for the discipline) channels is not specified,
// then the capacity of each channel will be equal to the quantity of data items that
// will be written to it.
//
// When specifying multiple input channel capacities, the value of the first
// one will be used.
func NewBenchmarker(handlersQuantity uint, inputCapacity ...uint) (*Benchmarker, error) {
	if handlersQuantity == 0 {
		return nil, ErrHandlersQuantityZero
	}

	bnch := &Benchmarker{
		handlersQuantity: handlersQuantity,
		inputCapacities:  append([]uint(nil), inputCapacity...),

		channels:      make(map[uint]chan uint),
		itemsQuantity: make(map[uint]uint),
		wg:            new(sync.WaitGroup),
	}

	return bnch, nil
}

// Returns the quantity of data handlers specified when the instance was created.
func (bnch *Benchmarker) HandlersQuantity() uint {
	return bnch.handlersQuantity
}

// Increases quantity of data items that will be written to the input (for the
// discipline) channel of the specified priority.
func (bnch *Benchmarker) AddItems(priority uint, quantity uint) {
	bnch.itemsQuantity[priority] += quantity
}

// Recreates the input (for the discipline) channels where data items will be written.
//
// Must be called before [Benchmarker.Play].
func (bnch *Benchmarker) Inputs() map[uint]<-chan uint {
	clear(bnch.channels)

	inputs := make(map[uint]<-chan uint, len(bnch.itemsQuantity))

	for priority := range bnch.itemsQuantity {
		channel := make(chan uint, bnch.inputCapacity(priority))

		bnch.channels[priority] = channel
		inputs[priority] = channel
	}

	return inputs
}

func (bnch *Benchmarker) inputCapacity(priority uint) uint {
	if len(bnch.inputCapacities) != 0 {
		return bnch.inputCapacities[0]
	}

	return bnch.itemsQuantity[priority]
}

// Writes data items to the input (for the discipline) channels and reads data items
// from the output channel of the discipline.
//
// Must be called after [Benchmarker.Inputs].
func (bnch *Benchmarker) Play(discipline Discipline[uint]) {
	bnch.handlers(discipline)
	bnch.writers()
	bnch.wg.Wait()
}

func (bnch *Benchmarker) writers() {
	for priority := range bnch.itemsQuantity {
		bnch.wg.Add(1)

		go bnch.writer(priority)
	}
}

func (bnch *Benchmarker) writer(priority uint) {
	defer bnch.wg.Done()
	defer close(bnch.channels[priority])

	for item := range bnch.itemsQuantity[priority] {
		bnch.channels[priority] <- item
	}
}

func (bnch *Benchmarker) handlers(discipline Discipline[uint]) {
	for range bnch.handlersQuantity {
		bnch.wg.Add(1)

		go bnch.handler(discipline)
	}
}

func (bnch *Benchmarker) handler(discipline Discipline[uint]) {
	defer bnch.wg.Done()

	for item := range discipline.Output() {
		discipline.Release(item.Priority)
	}
}
