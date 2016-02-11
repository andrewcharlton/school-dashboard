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

	"github.com/andrewcharlton/school-dashboard/analysis/student"
	"github.com/andrewcharlton/school-dashboard/analysis/subject"
)

// A Lookup collects object names and their ID numbers
// from the database. IDs are stored as strings to avoid the
// need for conversion when parsing query strings.
type Lookup struct {
	ID   string
	Name string
}

// A Database provides a wrapper to the school database.
type Database struct {

	// Database connection
	conn *sql.DB

	// Config Options
	Config Config

	// Cached items
	Levels   map[int]subject.Level
	Subjects map[int]*subject.Subject

	// National data
	Attainment8 map[string](map[string]student.Attainment8)

	// Cached items for menu lookups
	Dates       []Lookup
	Resultsets  []Lookup
	NatYears    []Lookup
	Ethnicities []string

	// OtherEths tells us whether an ethnic group should be
	// collapsed into the 'Other' category
	OtherEths map[string]bool

	// Maps for lookups
	dateMap      map[string]string
	resultsetMap map[string]string
	natYearMap   map[string]string

	// Prepared statements
	stmts map[string]*sql.Stmt
}
