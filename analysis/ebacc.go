package analysis

import "sort"

// An EBaccResult collects the results for each of the seperate
// EBacc areas: "E", "M", "S", "H", "L" and "All".
type EBaccResult struct {
	KS2     float64
	Results map[string]Result
}

// EBacc calculates whether or not a student was entered/achieved
// the EBacc.  Students need to achieve passes in English Language,
// Maths, 2 Sciences, a Humanity and a Language.
//
// Result keys used: EntB, AchB, Pts
func (s Student) EBacc() EBaccResult {

	entered := map[string]int{}
	passed := map[string]int{}
	points := map[string][]int{}
	results := map[string]Result{}

	for _, c := range s.Courses {
		entered[c.EBacc] += 1
		points[c.EBacc] = append(points[c.EBacc], c.Pts)
		if c.L2Pass {
			passed[c.EBacc] += 1
		}
	}

	// Sort points lists
	for _, pts := range points {
		pts = append(pts, 0) // Ensure that there are some values
		sort.Sort(sort.Reverse(sort.IntSlice(pts)))
	}

	entAll, passAll := true, true
	// Maths, Humanities, Languages
	for _, key := range []string{"M", "H", "L"} {
		entAll = entAll && (entered[key] >= 1)
		passAll = passAll && (passed[key] >= 1)
		results[key] = Result{EntB: entered[key] >= 1,
			AchB: passed[key] >= 1,
			Pts:  float64(points[key][0])}
	}

	// English - needs to be entered in both Language and Literature, and
	// achieve a pass in either one of them.
	switch {
	case !(entered["El"] >= 1 && entered["En"] >= 1):
		results["E"] = Result{EntB: false, AchB: false}
	case !(passed["El"] >= 1 || passed["En"] >= 1):
		results["E"] = Result{EntB: true, AchB: false}
	case points["En"][0] > points["El"][0]:
		results["E"] = Result{EntB: true, AchB: true, Pts: float64(points["En"][0])}
	default:
		results["E"] = Result{EntB: true, AchB: true, Pts: float64(points["El"][0])}
	}
	entAll = entAll && results["E"].EntB
	passAll = passAll && results["E"].AchB

	// Science - needs entries/passes in Core/Additional or Double Science
	// or any two of the triple sciences/computer science.
	switch {
	case entered["Sd"] >= 2:
		results["S"] = Result{EntB: true, AchB: passed["Sd"] >= 2,
			Pts: (float64(points["Sd"][0] + points["Sd"][1])) / float64(2)}
	case entered["St"] >= 2:
		results["S"] = Result{EntB: true, AchB: passed["St"] >= 2,
			Pts: (float64(points["St"][0] + points["St"][1])) / float64(2)}
	default:
		results["S"] = Result{EntB: false, AchB: false}
	}

	results["All"] = Result{EntB: entAll, AchB: passAll}
	return EBaccResult{KS2: s.KS2.APS, Results: results}
}

// EBacc calculates the number of students who were entered/achieved
// the EBacc.  Students need to achieve passes in English Language,
// Maths, 2 Sciences, a Humanity and a Language.
//
// Result keys used: EntN, AchN, EntP, AchP
func (g Group) EBacc() EBaccResult {

	cohort := float64(len(g.Students))

	// Initialise map
	keys := []string{"E", "M", "S", "H", "L", "All"}
	results := map[string]Result{}
	for _, key := range keys {
		results[key] = Result{}
	}

	if cohort == 0 {
		return EBaccResult{Results: results}
	}

	// Get student numbers, and points achieved by all students.
	points := map[string]float64{}
	for _, s := range g.Students {
		e := s.EBacc().Results
		for _, key := range keys {
			if e[key].EntB || key == "E" || key == "M" {
				pts := points[key]
				points[key] = pts + float64(e[key].Pts)
			}
			if e[key].EntB {
				r := results[key]
				r.EntN += 1
				results[key] = r
			}
			if e[key].AchB {
				r := results[key]
				r.AchN += 1
				results[key] = r
			}
		}
	}

	// English, Maths, Overall figures given as a percentage
	// of the cohort.
	for _, key := range []string{"E", "M", "All"} {
		r := results[key]
		r.EntP = float64(r.EntN) / cohort
		r.AchP = float64(r.AchN) / cohort
		r.Pts = points[key] / cohort
		results[key] = r
	}

	// Passrates in other subject given as a percentage of the
	// cohort actually studying that subject.
	for _, key := range []string{"S", "H", "L"} {
		r := results[key]
		r.EntP = float64(r.EntN) / cohort
		if r.EntN > 0 {
			r.AchP = float64(r.AchN) / float64(r.EntN)
			r.Pts = points[key] / float64(r.EntN)
		}
	}

	return EBaccResult{Results: results}
}
