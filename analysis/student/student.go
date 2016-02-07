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
	Gender     Gender
	Ethnicity  string
	SEN        SENInfo
	KS2        KS2Info
	Results    map[string]subject.Result
	Attendance AttendanceInfo

	// Calculated values
	basket *Basket // cached Progress8 basket

	// National values
	natYear       string
	natAtt8       Attainment8
	natAtt8Exists bool
}

// Name returns the student's name, formatted as
// Surname, Forename
func (s Student) Name() string {

	return s.Surname + ", " + s.Forename
}

// SetNationals allocates the national data to the student
func (s *Student) SetNationals(year string, att8 Attainment8, exists bool) {

	s.natYear = year
	s.natAtt8 = att8
	s.natAtt8Exists = exists
}

// Gender: coded 0 = Female, 1 = Male
type Gender int

// String converts the bool representation of the gender
// to "Male" or "Female"
func (g Gender) String() string {

	if g == 1 {
		return "Male"
	}
	return "Female"
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

// Score returns the student's ks2 score for a particular aspect
// En, Ma, Av are acceptable inputs.  Av used for any other input
func (ks2 KS2Info) Score(area string) string {

	switch {
	case area == "En" && ks2.En != "":
		return ks2.En
	case area == "En": // use reading score where there is no overall english
		return ks2.Re
	case area == "Ma":
		return ks2.Ma
	default:
		return ks2.Av
	}
}
