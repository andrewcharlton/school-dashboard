package subject

// A Grade
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