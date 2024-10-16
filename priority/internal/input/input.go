// Internal package with description of input channels parameters.
package input

// Input channel parameters.
type Input[Type any] struct {
	Channel <-chan Type
	Drained bool
}
