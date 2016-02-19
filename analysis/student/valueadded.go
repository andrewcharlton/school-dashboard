package student

import "fmt"

type VAScore struct {
	Expected float64
	Achieved float64
	Err      error
}

// Score calculates the actual VA score
func (va VAScore) Score() float64 {
	if va.Err != nil {
		return 0.0
	}
	return va.Achieved - va.Expected
}

// SubjectVa calculates the Value Added score for a student in a particular subject.
// VA is expressed in terms of grades above/below where we would expect students to be.
func (s Student) SubjectVA(subj string) VAScore {

	// Hacky ks3 - remove later
	if s.Year == 7 || s.Year == 8 || s.Year == 9 {
		return s.ks3VA(subj)
	}

	r, exists := s.Results[subj]
	if !exists {
		return VAScore{Err: fmt.Errorf("No Result for %v in %v found", s.Name(), subj)}
	}

	ks2 := s.KS2.Score(r.KS2Prior)
	if ks2 == "" {
		return VAScore{Err: fmt.Errorf("No %v KS2 score for %v", r.KS2Prior, s.Name())}
	}

	tm, exists := r.Subject.TMs[s.natYear]
	if !exists {
		return VAScore{Err: fmt.Errorf("No TM for %v in %v", subj, s.natYear)}
	}

	exp, err := tm.Expected(ks2)
	if err != nil {
		return VAScore{Err: err}
	}

	return VAScore{Expected: exp, Achieved: r.Att8}
}
