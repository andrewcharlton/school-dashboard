package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQL Driver
)

// A SchoolDB wraps an SQL object and provides additional
// functions to query the database.
type SchoolDB struct {
	DB *sql.DB
}

// Connect opens a connection to the database and returns
// a SchoolDB object.
func Connect(filename string) (SchoolDB, error) {

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return SchoolDB{}, nil
	}

	return SchoolDB{db}, nil

}

// Close terminates the connection to the database.
func (db SchoolDB) Close() error {

	err := db.DB.Close()
	return err
}
