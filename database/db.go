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

// A Database provides a wrapper for an sqlite database.
type Database struct {

	// Database connection
	conn  *sql.DB
	stmts map[string]*sql.Stmt

	// Cached object
	levels   map[int]subject.Level
	subjects map[int]subject.Subject
}

// ConnectDB opens a connection to the database
// and returns a Database object
func ConnectDB(filename string) (Database, error) {

	var db Database

	conn, err := sql.Open("sqlite3", filename)
	if err != nil {
		return db, err
	}
	db.conn = conn

	// Prepare SQL statements
	if err := db.prepareStatements(); err != nil {
		return db, err
	}

	// Load details of levels cache
	if err := db.loadLevels(); err != nil {
		return db, err
	}

	// Load subjects to cache
	if err := db.loadSubjects(); err != nil {
		return db, err
	}

	return db, nil
}

// Close terminates the connection to the
// database.
func (db Database) Close() error {

	err := db.conn.Close()
	return err
}

// loadLevels pulls in all of the levels for caching
func (db *Database) loadLevels() error {

	rows, err := db.conn.Query("SELECT id, level, is_gcse FROM levels")
	if err != nil {
		return err
	}
	defer rows.Close()

	db.levels = map[int]subject.Level{}
	for rows.Next() {
		var id int
		var l subject.Level
		err := rows.Scan(&id, &l.Lvl, &l.IsGCSE)
		if err != nil {
			return err
		}

		grades, err := db.loadGrades(id)
		if err != nil {
			return err
		}
		l.Gradeset = grades

		db.levels[id] = l
	}
	return nil
}

// loadGrades pulls in the list of grades at a particular level
func (db *Database) loadGrades(lvl int) (map[string]lvl.Grade, error) {

	rows, err := db.conn.Query(`SELECT grade, points, att8, l1_pass, l2_pass
								FROM grades
								WHERE level_id=?`, lvl)
	if err != nil {
		return map[string]lvl.Grade{}, err
	}
	defer rows.Close()

	grades := map[string]lvl.Grade{}
	for rows.Next() {
		var g lvl.Grade
		err := rows.Scan(&g.Grd, &g.Pts, &g.Att8, &g.L1Pass, &g.L2Pass)
		if err != nil {
			return map[string]lvl.Grade{}, err
		}
		grades[g.Grd] = g
	}
	return grades, nil
}

// loadSubjects pulls in the subject list for caching
func (db *Database) loadSubjects() error {

	rows, err := db.conn.Query(`SELECT id, subject, level_id, ebacc, tm, ks2_prior
								FROM subjects`)
	if err != nil {
		return err
	}
	defer rows.Close()

	db.subjects = map[int]*analysis.Subject{}
	for rows.Next() {
		var lvl int
		var s analysis.Subject
		err := rows.Scan(&s.SubjID, &s.Subj, &lvl, &s.EBacc, &s.TM, &s.KS2Prior)
		if err != nil {
			return err
		}
		s.Level = db.levels[lvl]
		db.subjects[s.SubjID] = &s
	}
	return nil
}
