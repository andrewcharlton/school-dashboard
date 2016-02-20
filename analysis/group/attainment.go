package group

import "github.com/andrewcharlton/school-dashboard/analysis/student"

// A Result contains the number of students that are eligible for, and
// the number who achieved a particular measure.
type Result struct {
	Entered  int
	Achieved int
	Percent  float64
}

// AveragePoints returns the average number of points achieved in
// a subject.
func (g Group) AveragePoints(subj string) float64 {

	total, cohort := 0, 0
	for _, s := range g.Students {
		r, exists := s.Results[subj]
		if !exists {
			continue
		}
		total += r.Pts
		cohort++
	}

	if cohort == 0 {
		return 0.0
	}
	return float64(total) / float64(cohort)
}

// Basics calculates the percentages of students in the group
// achieving a Level 2 Pass in both English and Maths,
func (g Group) Basics() Result {

	passes, entered := 0, 0
	for _, s := range g.Students {
		entered++
		if s.Basics() {
			passes++
		}
	}

	if entered == 0 {
		return Result{0, 0, 0.0}
	}
	return Result{entered, passes, float64(passes) / float64(entered)}
}

// An EBaccSummary contains the details of how many students in a cohort
// entered/achieved the various elements of the EBacc, and the associated
// percentages.
type EBaccSummary struct {
	Entered    int
	Achieved   int
	PctCohort  float64
	PctEntries float64
}

// EBaccArea provides a summary of the group's performance in one
// element of the EBacc
func (g Group) EBaccArea(area string) EBaccSummary {

	return g.ebaccSummary(area)
}

// EBacc provides a summary of the group's performance across the
// whole EBacc.
func (g Group) EBacc() EBaccSummary {

	return g.ebaccSummary("")
}

func (g Group) ebaccSummary(area string) EBaccSummary {

	if len(g.Students) == 0 {
		return EBaccSummary{0, 0, 0.0, 0.0}
	}

	entered, achieved := 0, 0
	for _, s := range g.Students {
		var r student.EBaccResult
		switch area {
		case "":
			r = s.EBacc()
		default:
			r = s.EBaccArea(area)
		}
		if r.Entered {
			entered++
		}
		if r.Achieved {
			achieved++
		}
	}

	if entered == 0 {
		return EBaccSummary{0, 0, 0.0, 0.0}
	}

	return EBaccSummary{entered, achieved,
		float64(achieved) / float64(len(g.Students)),
		float64(achieved) / float64(entered),
	}
}
