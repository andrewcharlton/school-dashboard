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

// A List allows subjects to be sorted by name and then level.
type List []Subject

func (l List) Len() int      { return len(l) }
func (l List) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l List) Less(i, j int) bool {
	switch {
	case l[i].Subj < l[j].Subj:
		return true
	case l[j].Subj < l[i].Subj:
		return false
	case l[i].Lvl < l[j].Lvl:
		return true
	default:
		return false
	}
}
