package group

import "github.com/andrewcharlton/school-dashboard/analysis/student"

var (
	// PP = Disadvantaged students
	PP = func(s student.Student) bool { return s.PP }

	// NonPP = Non-Disadvantaged students
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

	// PA = Persistent absentees (less than 90% attendance)
	PA = func(s student.Student) bool { return s.Attendance.Latest() < 0.9 }

	// NonPA = Students with greater than 90% attendance
	NonPA = func(s student.Student) bool { return s.Attendance.Latest() >= 0.9 }

	NonSEN = func(s student.Student) bool { return s.SEN.Status == "N" }

	SEN = func(s student.Student) bool { return s.SEN.Status == "K" || s.SEN.Status == "S" }

	Statement = func(s student.Student) bool { return s.SEN.Status == "S" }
)

// Year returns a subgroup filter for a certain yeargroup
func Year(year int) func(s student.Student) bool {
	return func(s student.Student) bool { return s.Year == year }
}

// Form returns a subgroup filter for a certain form group
func Form(form string) func(s student.Student) bool {
	return func(s student.Student) bool { return s.Form == form }
}

// Studying returns a subgroup filter for students studying a particular subject.
func Studying(subj string, subjID int) func(s student.Student) bool {
	return func(s student.Student) bool {
		r, exists := s.Results[subj]
		if !exists {
			return false
		}
		return r.SubjID == subjID
	}
}

// Class returns a subgroup filter for a certain subject/class combination
func Class(subject, class string) func(s student.Student) bool {
	return func(s student.Student) bool {
		r, exists := s.Results[subject]
		if !exists {
			return false
		}
		return r.Class == class
	}
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
