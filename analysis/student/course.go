package student

import "github.com/andrewcharlton/school-dashboard/analysis/lvl"

// A Course brings together a subject and the grade
// achieved in that subject.
type Course struct {
	*Subject
	*lvl.Grade
	Effort  int
	Class   string
	Teacher string
}
