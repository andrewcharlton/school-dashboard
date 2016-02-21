// Package subject pulls together all of the
package subject

// A Subject contains subject details.
type Subject struct {

	// Name of the subject
	Subj string

	// Short identifier
	Code string

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

// A SubjectList allows subjects to be sorted by name and then level.
type SubjectList []Subject

func (s SubjectList) Len() int      { return len(s) }
func (s SubjectList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SubjectList) Less(i, j int) bool {
	switch {
	case s[i].Subj < s[j].Subj:
		return true
	case s[j].Subj < s[i].Subj:
		return false
	case s[i].Lvl < s[j].Lvl:
		return true
	default:
		return false
	}
}
