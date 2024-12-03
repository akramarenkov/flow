package inspect

// Options of inspecting.
type Opts struct {
	// Quantity or maximum quantity of data handlers that will be passed to the divider
	// when inspecting
	Quantity uint
	// Priority list on which the divider will be inspected
	Priorities []uint
}

// Inspection result.
type Result struct {
	// Inspection conclusion. Filled in case of an error occurs in filling the
	// distribution by the inspected divider, in case of an error occurs in
	// calculation of the parameters of the obtained distribution or in case
	// discrepancy between the distribution returned by the inspected divider and
	// the expected distribution
	Conclusion error
	// Error returned by the inspected divider or other function participating in the
	// inspection
	Err error
	// Quantity of data handlers for which a negative inspection conclusion was obtained
	Quantity uint
	// Priority list for which a negative inspection conclusion was obtained
	Priorities []uint
}
