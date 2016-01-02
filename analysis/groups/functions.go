package groups

import (
	"fmt"

	"github.com/andrewcharlton/school-dashboard/analysis"
)

// A ScoreFunc produces a score from a student.
// The name is used as a header when producing summary tables
type ScoreFunc struct {
	Name string
	Func func(analysis.Student) Score
}

// AttendancePct returns the % Attendance recorded for the year.
var AttendancePct = ScoreFunc{
	Name: "% Att",
	Func: func(s analysis.Student) Score {
		return Score{s.Attendance.Latest(), false}
	},
}

// AttendanceUnder returns a ScoreFunc which calculates the
// % of students currently with an attendance under a certain
// %.
func AttendanceUnder(att float64) ScoreFunc {
	return ScoreFunc{
		Name: fmt.Sprintf("Att Under %.0f%%", att),
		Func: func(s analysis.Student) Score {
			if s.Attendance.Latest() < att {
				return Score{100.0, false}
			}
			return Score{0.0, false}
		},
	}
}
