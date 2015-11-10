package analysis

// AttendanceInfo collects a student's latest attendance
// data together
type AttendanceInfo struct {
	Possible     int
	Absences     int
	Unauthorised int
	Sessions     [10]int
}

// Latest calculates the most recent attendance figure for
// a student, expressed as a percentage.
func (att AttendanceInfo) Latest() float64 {

	if att.Possible == 0 {
		return float64(0)
	}

	return float64(100) * float64(att.Possible-att.Absences) / float64(att.Possible)
}

var AttendanceSessions = map[int]string{
	0: "Mon AM",
	1: "Mon PM",
	2: "Tue AM",
	3: "Tue PM",
	4: "Wed AM",
	5: "Wed PM",
	6: "Thu AM",
	7: "Thu PM",
	8: "Fri AM",
	9: "Fri PM",
}
