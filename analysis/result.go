package analysis

// A Result is used to wrap the results from any calculations.
// This is used so that methods can be called from templates.
// Not all fields will be used in all calculations.
type Result struct {

	// Whether students were entered/eligible for a result
	// and whether they achieved it.
	// E.g whether they were entered for EBacc and/or achieved it.
	EntB bool
	AchB bool

	// How many eligible entries students had, and how many
	// they achieved.
	// Example: How many GCSEs students took, and how many they
	// passed.
	EntN int
	AchN int

	// What percentage of
	EntP float64
	AchP float64

	// How many points a student achieved in a measure, and
	// how many they were expected to achieve.
	// Example: Points in the Progress 8 Basket, and the national
	// figure for similar students.
	Pts float64
	Exp float64

	// Any errors thrown up in the calculation.
	Error error
}
