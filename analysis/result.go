package analysis

// A Result is used to wrap the results from any calculations.
// Not all fields will be used in all calculations.
type Result struct {
	Entered  bool
	Achieved bool
	Expected float64
	Points   float64
	Error    error
}
