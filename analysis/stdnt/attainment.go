package stdnt

// Basics measures whether a student has achieved a level 2 pass
// in both English and Maths.
func (s Student) Basics() (bool, bool) {

	eng, maths := false, false
	for _, c := range s.Courses {
		switch c.EBacc {
		case "En":
			if c.L2Pass {
				eng = true
			}
		case "El":
			if c.L2Pass {
				eng = true
			}
		case "M":
			if c.L2Pass {
				maths = true
			}
		}
	}

	return true, eng && maths
}

func (s Student) EBacc() (bool, bool) {

	eEng, aEng := s.EBaccEng()
	eMat, aMat := s.EBaccMaths()
	eSci, aSci := s.EBaccSci()
	eHum, aHum := s.EBaccHum()
	eLan, aLan := s.EBaccLang()
	entered := eEng && eMat && eSci && eHum && eLan
	achieved := aEng && aMat && aSci && aHum && aLan
	return entered, achieved

}

// EBaccEng calculates whether a student has achieved a level
// 2 pass in the English section of the EBacc
func (s Student) EBaccEng() (bool, bool) {

	entLang, entLit := false, false
	achLang, achLit := false, false
	for _, c := range s.Courses {
		switch c.EBacc {
		case "E":
			entLang, entLit = true, true
			if c.L2Pass {
				achLang = true
				achLit = true
			}
		case "En":
			entLang = true
			if c.L2Pass {
				achLang = true
			}
		case "El":
			entLit = true
			if c.L2Pass {
				achLit = true
			}
		}
	}
	return entLang && entLit, achLang && achLit
}

// EBaccSci calculates whether or not a student was entered for/
// achieved two Science qualifications
func (s Student) EBaccSci() (bool, bool) {

	entries, passes := 0, 0
	for _, c := range s.Courses {
		if c.EBacc == "S" {
			entries++
			if c.L2Pass {
				passes++
			}
		}
	}
	return entries >= 2, passes >= 2
}

// EBaccMaths calculates whether or not a student was entered for/
// achieved a Mathematics qualification.
func (s Student) EBaccMaths() (bool, bool) {
	return s.ebaccArea("M")
}

// EBaccLang calculates whether or not a student was entered for/
// achieved a Language qualification.
func (s Student) EBaccLang() (bool, bool) {
	return s.ebaccArea("L")
}

// EBaccHum calculates whether or not a student was entered for/
// achieved a Humanities qualification.
func (s Student) EBaccHum() (bool, bool) {
	return s.ebaccArea("H")
}

// ebaccArea helper function
func (s Student) ebaccArea(area string) (bool, bool) {

	ent, ach := false, false
	for _, c := range s.Courses {
		switch c.EBacc {
		case area:
			ent = true
			if c.L2Pass {
				ach = true
			}
		}
	}
	return ent, ach
}
