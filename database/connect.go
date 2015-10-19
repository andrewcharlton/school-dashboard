package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQL Driver
)

// A SchoolDB wraps an SQL object and provides additional
// functions to query the database.
type SchoolDB struct {
	db      *sql.DB
	queries map[string]*sql.Stmt
}

// Connect opens a connection to the database and returns
// a SchoolDB object.
func Connect(filename string) (SchoolDB, error) {

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return SchoolDB{}, nil
	}

	s := SchoolDB{db, map[string]*sql.Stmt{}}
	err = s.prepQueries()
	if err != nil {
		return SchoolDB{}, nil
	}

	return s, nil
}

// Close terminates the connection to the database.
func (db SchoolDB) Close() error {

	err := db.closeQueries()
	if err != nil {
		return err
	}

	err = db.DB.Close()
	return err
}
