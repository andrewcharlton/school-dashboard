package analysis

import "github.com/andrewcharlton/school-dashboard/level"

// A Subject contains subject details.
type Subject struct {

	// Name of the subject
	Subj string

	// The Level of the subject
	*level.Level

	// Which EBacc Basket the Subject falls in:
	// M: Maths
	// En: English Language, El: English Lit
	// S: Science
	// H: Humanities
	// L: Languages
	// Empty for non-EBacc subjects
	EBacc string

	// Which Transition Matrix to use for the subject.
	TM string

	// Which KS2 score the subject TM is based on:
	// En, Ma, Av
	KS2Prior string
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
