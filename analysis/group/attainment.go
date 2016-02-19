package group

// A Result has
type Result struct {
	Ent int
	Ach int
	Pct float64
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
		return Result{Ent: 0, Ach: 0, Pct: 0.0}
	}
	return Result{Ent: entered, Ach: passes, Pct: float64(passes) / float64(entered)}
}
