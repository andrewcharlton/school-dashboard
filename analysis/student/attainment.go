package student

import "fmt"

// SubjectGrade provides a lookup for the grade in a particular subject.
// Returns "" if no grade is found.
func (s Student) SubjectGrade(subj string) string {

	r, exists := s.Results[subj]
	if !exists {
		return ""
	}
	return r.Grd
}

// APS returns the average point score achieved by a student in his/her courses.
func (s Student) APS() float64 {

	total, number := 0, 0
	for _, r := range s.Results {
		if r.Grd != "" {
			total += r.Pts
			number++
		}
	}
	if number == 0 {
		return 0.0
	}
	return float64(total) / float64(number)
}

// Basics measures whether a student has achieved a level 2 pass
// in both English and Maths.
func (s Student) Basics() bool {

	eng, maths := false, false
	for _, r := range s.Results {
		switch r.EBacc {
		case "En", "El", "E":
			if r.L2Pass {
				eng = true
			}
		case "M":
			if r.L2Pass {
				maths = true
			}
		}
	}
	return eng && maths
}

// An EBaccResult wraps the results from the EBacc calculations so
// they can be used directly in a template.
type EBaccResult struct {
	Entered  bool
	Achieved bool
	Results  []string
}

// EBacc calculates whether a student is eligible for the EBacc
// and whether or not they have achieved it.
func (s Student) EBacc() EBaccResult {

	entered, achieved := true, true
	for _, a := range []string{"E", "M", "S", "H", "L"} {
		r := s.EBaccArea(a)
		entered = entered && r.Entered
		achieved = achieved && r.Achieved
	}

	return EBaccResult{entered, achieved, []string{}}
}

// EBaccEntries calculates how many ebacc areas a student has entered.
func (s Student) EBaccEntries() int {

	entries := 0
	for _, area := range []string{"E", "M", "S", "H", "L"} {
		eb := s.EBaccArea(area)
		if eb.Entered {
			entries++
		}
	}
	return entries
}

// EBaccPasses calculates how many areas a student has achieved an
// EBacc pass.
func (s Student) EBaccPasses() int {

	passes := 0
	for _, area := range []string{"E", "M", "S", "H", "L"} {
		eb := s.EBaccArea(area)
		if eb.Achieved {
			passes++
		}
	}
	return passes
}

// EBaccArea calculates whether a student was entered and/or achieved
// a pass in the relevant section of the EBacc. Valid values for area are:
// * E: English
// * M: Maths
// * S: Science
// * H: Humanities
// * L: Languages
func (s Student) EBaccArea(area string) EBaccResult {

	switch area {
	case "E":
		return s.ebaccEng()
	case "S":
		return s.ebaccSci()
	default:
		return s.ebaccArea(area)
	}
}

// EBaccEng calculates whether a student has achieved a level
// 2 pass in the English section of the EBacc
func (s Student) ebaccEng() EBaccResult {

	entLang, entLit := false, false
	achLang, achLit := false, false
	results := []string{}
	for _, r := range s.Results {
		switch r.EBacc {
		case "E":
			results = append(results, fmt.Sprintf("%v (%v)", r.Subj, r.Grd))
			entLang, entLit = true, true
			if r.L2Pass {
				achLang = true
				achLit = true
			}
		case "En":
			results = append(results, fmt.Sprintf("%v (%v)", r.Subj, r.Grd))
			entLang = true
			if r.L2Pass {
				achLang = true
			}
		case "El":
			results = append(results, fmt.Sprintf("%v (%v)", r.Subj, r.Grd))
			entLit = true
			if r.L2Pass {
				achLit = true
			}
		}
	}
	return EBaccResult{entLang && entLit, achLang || achLit, results}
}

// EBaccSci calculates whether or not a student was entered for/
// achieved two Science qualifications
func (s Student) ebaccSci() EBaccResult {

	entries, passes := 0, 0
	results := []string{}
	for _, r := range s.Results {
		if r.EBacc == "S" {
			results = append(results, fmt.Sprintf("%v (%v)", r.Subj, r.Grd))
			entries++
			if r.L2Pass {
				passes++
			}
		}
	}
	return EBaccResult{entries >= 2, passes >= 2, results}
}

// ebaccArea helper function
func (s Student) ebaccArea(area string) EBaccResult {

	ent, ach := false, false
	results := []string{}
	for _, r := range s.Results {
		if r.EBacc == area {
			results = append(results, fmt.Sprintf("%v (%v)", r.Subj, r.Grd))
			ent = true
			if r.L2Pass {
				ach = true
			}
		}
	}
	return EBaccResult{ent, ach, results}
}
