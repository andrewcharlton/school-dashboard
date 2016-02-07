// Package subject pulls together all of the
package subject

// A Subject contains subject details.
type Subject struct {

	// Name of the subject
	Subj string

	// Unique ID of the subject
	SubjID int

	// The Level of the subject
	Level

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

	// tms for the subject
	TMs map[string]TransitionMatrix
}
