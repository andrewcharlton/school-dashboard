// Package database implements a database access layer
// for the school-dashboard application.
//
// It wraps an SQL database connection, prepares
// statements to query the database and then
// provides methods to return the results of any
// queries.
package database

import (
	"database/sql"

	"github.com/andrewcharlton/school-dashboard/analysis/subject"
)

type Database struct {

	// Database connection
	conn  *sql.DB
	stmts map[string]*sql.Stmt

	// Cached object
	subjects map[int]subject.Subject
	levels   map[int]subject.Level
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
	NatYear   string
	Year      string
}

func (cfg Config) DefaultFilter() Filter {

	return Filter{Date: cfg.Date,
		Resultset: cfg.Resultset,
		NatYear:   cfg.NatYear,
		Year:      cfg.Year,
	}
}
