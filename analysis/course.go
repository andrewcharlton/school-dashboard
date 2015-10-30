package analysis

// A Subject contains subject details.
type Subject struct {

	// Name of the subject
	Subj string

	// The Level of the subject
	*Level

	// Which EBacc Basket the Subject falls in:
	// M: Maths
	// En: English Language, El: English Lit
	// S: Science
	// H: Humanities
	// L: Languages
	// Empty for non-EBacc subjects
	EBacc string

	// Which KS2 score the subject TM is based on:
	// En, Ma, Av
	KS2Prior string
}

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

// A Course brings together a subject and the grade
// achieved in that subject.
type Course struct {
	*Subject
	*Grade
	Effort  int
	Class   string
	Teacher string
}
