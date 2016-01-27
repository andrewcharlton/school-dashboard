package database

// A Filter holds all the fields available to select a group
// of students on.
type Filter struct {
	// Effective date to use for membership of the group.
	// Holds the id number, but held as string for easier
	// parsing from query strings.
	Date string
	// Which set of assessment results to use.
	// Holds id number as above.
	Resultset string
	// Which year to use for national data comparisons.
	NatYear string
	// Which yeargroup to include in the group
	Year string
	// Pupil premium: "", "1", or "0" for Any/True/False
	PP string
	// EAL students: "", "1", or "0" for Any/True/False
	EAL string
	// Gender: "", "1", "0" for Any/Male/Female
	Gender string
	// SEN types to include - Empty for any
	SEN []string
	// Which Ethnic groups to include, empty for any
	Ethnicities []string
	// Which KS2 bands to include, empty for any
	KS2Bands []string
}
