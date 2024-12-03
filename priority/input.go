package priority

// Discipline input channel descriptor.
type input[Type any] struct {
	Channel <-chan Type
	Closed  bool
}
