package group

import "github.com/andrewcharlton/school-dashboard/analysis/student"

var (
	// Disadvantaged students
	PP = func(s student.Student) bool { return s.PP }

	// Non-Disadvantaged students
	NonPP = func(s student.Student) bool { return !s.PP }

	// Male students
	Male = func(s student.Student) bool { return s.Gender == 1 }

	// Female students
	Female = func(s student.Student) bool { return s.Gender == 0 }

	// High attaining students at KS2
	High = func(s student.Student) bool { return s.KS2.Band == "High" }

	// Middle attaining students at KS2
	Middle = func(s student.Student) bool { return s.KS2.Band == "Middle" }

	// Low attainin students at KS2
	Low = func(s student.Student) bool { return s.KS2.Band == "Low" }

	// Persistent absentees (less than 90% attendance)
	PA = func(s student.Student) bool { return s.Attendance.Latest() < 0.9 }

	// Students with greater than 90% attendance
	NonPA = func(s student.Student) bool { return s.Attendance.Latest() >= 0.9 }
)

// Year returns a subgroup filter for a certain yeargroup
func Year(year int) func(s student.Student) bool {

	return func(s student.Student) bool { return s.Year == year }

}

// SubGroup produces a Group with a subset of the original students, as
// determined by the filter functions.
func (g Group) SubGroup(filters ...func(s student.Student) bool) Group {

	students := []student.Student{}
	for _, s := range g.Students {
		include := true
		for _, f := range filters {
			if !f(s) {
				include = false
				break
			}
		}
		if include {
			students = append(students, s)
		}
	}

	return Group{Students: students}
}