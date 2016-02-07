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

// A Database provides a wrapper to the school database.
type Database struct {

	// Database connection
	conn *sql.DB

	// Cached items
	levels   map[int]subject.Level
	subjects map[int]*subject.Subject
}
