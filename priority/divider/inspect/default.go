package inspect

// Returns default set of inspecting options.
func DefaultSet() []Opts {
	set := []Opts{
		{
			//nolint:mnd // For a specified quantity and for the associated priority
			// list the error of the distribution in integers does not exceed 10% of
			// the real value for the dividers [divider.Fair] and [divider.Rate]
			Quantity:   1000,
			Priorities: []uint{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
	}

	return set
}
