// Package db implements a database access layer
// for the school-dashboard application.
//
// It wraps an SQL database connection, prepares
// statements to query the database and then
// provides methods to return the results of any
// queries.
package db

// A Database provides a wrapper to the database
// connection and provides methods to query it.
type Database interface {
	Close() error
	Dates() ([]Lookup, error)
	Resultsets() ([]Lookup, error)
	Ethnicities() ([]Ethnicity, error)
	Group(f Filter) ([]string, error)
}

// A Lookup holds an ID/Name pair for an item in the database.
type Lookup struct {
	ID   string
	Name string
}

// An Ethnicity holds the name of an Ethnic group, as well as the
// number of students in that group.
type Ethnicity struct {
	Name  string
	Count int
}

// A Filter holds all the fields available to select a group
// of students on.
type Filter struct {
	// Effective date to use for membership of the group
	Date string
	// Which set of assessment results to use
	Resultset string
	// Which yeargroups to include in the group
	Years []string
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

// A StudentFilter holds all the fields necessary to uniquely
// identify a student and the grades etc.
type StudentFilter struct {
	UPN       string
	Date      string
	Resultset string
}
