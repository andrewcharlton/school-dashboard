// Package db implements a database access layer
// for the school-dashboard application.
//
// It wraps an SQL database connection, prepares
// statements to query the database and then
// provides methods to return the results of any
// queries.
package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQL Driver.
)

// A Database provides a wrapper to the database
// connection and provides methods to query it.
type Database struct {
	db *sql.DB
}

// Connect opens a connection to the database
// and returns a Database object
func Connect(filename string) (Database, error) {

	conn, err := sql.Open("sqlite3", filename)
	if err != nil {
		return SchoolDB{}, nil
	}

	db := Database{conn}
	return db, nil
}

// Close terminates the connection to the
// database.
func (db Database) Close() error {

	err = db.DB.Close()
	return err
}
