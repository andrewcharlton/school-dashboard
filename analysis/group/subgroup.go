package group

import "github.com/andrewcharlton/school-dashboard/analysis/student"

var (
	// Year 7 students
	Year7 = func(s student.Student) bool { return s.Year == 7 }

	// Year 8 students
	Year8 = func(s student.Student) bool { return s.Year == 8 }

	// Year 9 students
	Year9 = func(s student.Student) bool { return s.Year == 9 }

	// Year 10 students
	Year10 = func(s student.Student) bool { return s.Year == 10 }

	// Year 11 students
	Year11 = func(s student.Student) bool { return s.Year == 11 }

	// Disadvantaged students
	PP = func(s student.Student) bool { return s.PP }

	// Non-Disadvantaged students
	NonPP = func(s student.Student) bool { return !s.PP }

	// Male students
	Boys = func(s student.Student) bool { return s.Gender == "Male" }

	// Female students
	Girls = func(s student.Student) bool { return s.Gender == "Female" }

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

// SubGroup produces a Group with a subset of the original students, as
// determined by the filter functions.
func (g Group) SubGroup(filters ...func(s student.Student) bool) Group {

	students := []student.Student{}
	for _, s := range g.Students {
		include = true
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
