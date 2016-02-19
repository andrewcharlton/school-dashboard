package student

// SubjectGrade provides a lookup for the grade in a particular subject.
// Returns "" if no grade is found.
func (s Student) SubjectGrade(subj string) string {

	r, exists := s.Results[subj]
	if !exists {
		return ""
	}
	return r.Grd
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
}

// EBacc calculates whether a student is eligible for the EBacc
// and whether or not they have achieved it.
func (s Student) EBacc() EBaccResult {

	e := s.EBaccEng()
	m := s.EBaccMaths()
	c := s.EBaccSci()
	h := s.EBaccHum()
	l := s.EBaccLang()
	entered := e.Entered && m.Entered && c.Entered && h.Entered && l.Entered
	achieved := e.Achieved && m.Achieved && c.Achieved && h.Achieved && l.Achieved
	return EBaccResult{entered, achieved}

}

// EBaccEng calculates whether a student has achieved a level
// 2 pass in the English section of the EBacc
func (s Student) EBaccEng() EBaccResult {

	entLang, entLit := false, false
	achLang, achLit := false, false
	for _, r := range s.Results {
		switch r.EBacc {
		case "E":
			entLang, entLit = true, true
			if r.L2Pass {
				achLang = true
				achLit = true
			}
		case "En":
			entLang = true
			if r.L2Pass {
				achLang = true
			}
		case "El":
			entLit = true
			if r.L2Pass {
				achLit = true
			}
		}
	}
	return EBaccResult{entLang && entLit, achLang && achLit}
}

// EBaccSci calculates whether or not a student was entered for/
// achieved two Science qualifications
func (s Student) EBaccSci() EBaccResult {

	entries, passes := 0, 0
	for _, r := range s.Results {
		if r.EBacc == "S" {
			entries++
			if r.L2Pass {
				passes++
			}
		}
	}
	return EBaccResult{entries >= 2, passes >= 2}
}

// EBaccMaths calculates whether or not a student was entered for/
// achieved a Mathematics qualification.
func (s Student) EBaccMaths() EBaccResult {
	return s.ebaccArea("M")
}

// EBaccLang calculates whether or not a student was entered for/
// achieved a Language qualification.
func (s Student) EBaccLang() EBaccResult {
	return s.ebaccArea("L")
}

// EBaccHum calculates whether or not a student was entered for/
// achieved a Humanities qualification.
func (s Student) EBaccHum() EBaccResult {
	return s.ebaccArea("H")
}

// ebaccArea helper function
func (s Student) ebaccArea(area string) EBaccResult {

	ent, ach := false, false
	for _, r := range s.Results {
		switch r.EBacc {
		case area:
			ent = true
			if r.L2Pass {
				ach = true
			}
		}
	}
	return EBaccResult{ent, ach}
}
