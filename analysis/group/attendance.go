package group

// AttendanceSummary contains summary data for a group's attendance.
type AttendanceSummary struct {
	Cohort       int
	Week         float64
	Possible     int
	Absences     int
	Unauthorised int
	PAs          int
}

// PercentAttendance calculates the overall percentage attendance for the group.
func (a AttendanceSummary) PercentAttendance() float64 {

	if a.Possible == 0 {
		return 0.0
	}
	return 1.0 - float64(a.Absences)/float64(a.Possible)
}

// PercentAuthorised calculates the overall percentage of authorised absence
// for the group.
func (a AttendanceSummary) PercentAuthorised() float64 {

	if a.Possible == 0 {
		return 0.0
	}
	return float64(a.Absences-a.Unauthorised) / float64(a.Possible)
}

// PercentUnauthorised calculates the overall percentage of authorised absence
// for the group.
func (a AttendanceSummary) PercentUnauthorised() float64 {

	if a.Possible == 0 {
		return 0.0
	}
	return float64(a.Unauthorised) / float64(a.Possible)
}

// PercentPA calculates the percentage of the cohort with attendance under 90%.
func (a AttendanceSummary) PercentPA() float64 {

	if a.Cohort == 0 {
		return 0.0
	}
	return float64(a.PAs) / float64(a.Cohort)
}

// Attendance calculates the summary attendance figures for the group.
func (g Group) Attendance() AttendanceSummary {

	att := AttendanceSummary{Cohort: len(g.Students)}
	for _, s := range g.Students {
		att.Week += s.Attendance.Week
		att.Possible += s.Attendance.Possible
		att.Absences += s.Attendance.Absences
		att.Unauthorised += s.Attendance.Unauthorised
		att.Possible += s.Attendance.Possible
		if s.Attendance.Latest() < 0.9 {
			att.PAs++
		}
	}
	att.Week /= float64(att.Cohort)
	return att
}
