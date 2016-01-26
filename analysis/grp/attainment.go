package grp

// A Result has
type Result struct {
	Ent int
	Ach int
	Pct float64
}

// Basics calculates the percentages of students in the group
// achieving a Level 2 Pass in both English and Maths,
func (g Group) Basics() Result {

	passes, entered := 0, 0
	for _, s := range g.Students {
		ent, ach := s.Basics()
		if ent {
			entered++
		}
		if ach {
			passes++
		}
	}

	if entered == 0 {
		return Result{Ent: 0, Ach: 0, Pct: 0.0}
	}
	return Result{Ent: entered, Ach: passes, Pct: float64(passes) / float64(entered)}
}
