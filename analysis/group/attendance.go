package group

// AttendanceSummary contains summary data for a group's attendance.
type AttendanceSummary struct {
	Cohort       int
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
	return 100.0 - 100.0*float64(a.Absences)/float64(a.Possible)
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

	att := AttendanceSummary{}
	for _, s := range g.Students {
		att.Cohort++
		att.Possible += s.Attendance.Possible
		att.Absences += s.Attendance.Absences
		att.Unauthorised += s.Attendance.Unauthorised
		att.Possible += s.Attendance.Possible
		if s.Attendance.Latest() < 0.9 {
			att.PAs++
		}
	}
	return att
}
