package analysis

// A Student holds the relevant data for a single
// student, at a single point in time.  It holds
// all of their personal data, as well as the details
// of all the courses they are studying and the
// grades they are achieving.
type Student struct {
	UPN       string
	Surname   string
	Forename  string
	Year      int
	Form      string
	PP        bool
	EAL       bool
	Gender    string
	SEN       SenInfo
	Ethnicity string
	KS2       KS2Info
}

// SENInfo collects all of a student's SEN details
// together.
type SENInfo struct {
	Exists     bool
	Status     string
	Info       string
	Strategies string
}

// KS2Info collects all of a student's KS2 scores
// together.
type KS2Info struct {
	Exists bool
	APS    float64
	Band   string
	Ma     string
	En     string
	Av     string
}
