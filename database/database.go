// Package database implements a database access layer
// for the school-dashboard application.
//
// It wraps an SQL database connection, prepares
// statements to query the database and then
// provides methods to return the results of any
// queries.
package database

import "github.com/andrewcharlton/school-dashboard/analysis"

// A Database provides a wrapper to the database
// connection and provides methods to query it.
type Database interface {

	// Close the database connection
	Close() error

	// Config returns a populated config object from the
	// database.
	Config() (Config, error)

	// Dates returns a sorted list of all effective dates
	// in the database that are marked to be listed.
	Dates() ([]Lookup, error)

	// Ethnicities returns a list of all the distinct
	// ethnicities present in the database, and the frequency
	// that each appears with.  Only students who are present
	// in listed dates are counted.  List should be sorted
	// with largest groups first.
	Ethnicities() ([]Ethnicity, error)

	// Resultsets returns a sorted list of all resultsets
	// in the database.
	Resultsets() ([]Lookup, error)

	// Student returns a student object with details relevant
	// to the given filter.
	Student(f StudentFilter) (analysis.Student, error)

	// Group returns a list of students who satisfy the
	// criteria specified in the filter
	Group(f Filter) (analysis.Group, error)

	// Search returns a list of students from a search query.
	Search(name, date string) ([]StudentLookup, error)
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

// A StudentLookup holds the basic details of a student,
// returned from a search
type StudentLookup struct {
	UPN  string
	Name string
	Year string
	Form string
}

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

// A Config holds all of the
type Config struct {

	// Name of the school
	School string

	// URN/LEA numbers (as found on performance tables)
	URN string
	LEA string

	// Default filter options
	Date      string
	Resultset string
	Years     []string
}

func (cfg Config) DefaultFilter() Filter {

	return Filter{Date: cfg.Date,
		Resultset: cfg.Resultset,
		Years:     cfg.Years,
	}
}
