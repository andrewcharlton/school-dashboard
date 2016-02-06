package student

import "github.com/andrewcharlton/school-dashboard/analysis/subject"

// A Student holds the relevant data for a single
// student, at a single point in time.  It holds
// all of their personal data, as well as the details
// of all the courses they are studying and the
// grades they are achieving.
type Student struct {
	UPN        string
	Surname    string
	Forename   string
	Year       int
	Form       string
	PP         bool
	EAL        bool
	Gender     string
	Ethnicity  string
	SEN        SENInfo
	KS2        KS2Info
	Results    map[string]subject.Result
	Attendance AttendanceInfo

	basket  *Basket     // cached Progress8 basket
	natAtt8 Attainment8 // National Attainment8 scores

}

// Name returns the student's name, formatted as
// Surname, Forename
func (s Student) Name() string {

	return s.Surname + ", " + s.Forename
}

// SENInfo collects all of a student's SEN details
// together.
type SENInfo struct {
	Status     string
	Need       string
	Info       string
	Strategies string
	Access     string
	IEP        bool
}

// KS2Info collects all of a student's KS2 scores
// together.
type KS2Info struct {
	Exists bool
	APS    float64
	Band   string
	En     string
	Ma     string
	Av     string
	Re     string
	Wr     string
	GPS    string
}
