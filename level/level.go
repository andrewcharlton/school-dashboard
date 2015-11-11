// Package level extracts the level struct so it
// can be used in both analysis and national packages.
package level

// A Level brings together the details and grades
// of a particular qualification level.
type Level struct {

	// Name
	Lvl string

	// Does the qualification count as a GCSE
	IsGCSE bool

	// Possible grades achievable at that level
	Gradeset map[string]*Grade
}

// A Grade contains the various points values for a
// grade at a particular level.
type Grade struct {

	// Name of the grade
	Grd string

	// QCA points
	Pts int

	// Attainment 8 points
	Att8 float64

	// Whether the grade counts as a pass at levels 1/2
	L1Pass bool
	L2Pass bool
}
