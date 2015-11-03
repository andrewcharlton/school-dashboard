package analysis

import "errors"

// Basics measures whether a student has achieved a level 2 pass
// in both English and Maths.
func (s Student) Basics() Result {

	eng, maths := false, false
	for _, c := range s.Courses {
		switch c.EBacc {
		case "En":
			if c.L2Pass {
				eng = true
			}
		case "M":
			if c.L2Pass {
				maths = true
			}
		}
	}

	return Results{Achieved: eng && maths}
}

// Basics calculates the percentages of students in the group
// achieving a Level 2 Pass in both English and Maths,
func (g Group) Basics() Results {

	passes, entered := 0, 0
	for _, s := range g.Students {
		if s.Basics().Achieved {
			passes += 1
		}
		entered += 1
	}

	if entered == 0 {
		return Result{Points: float64(0), Error: errors.New("No students in group")}
	}
	return Result{Points: float64(passes) / float64(entered), Error: nil}
}
