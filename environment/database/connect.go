package database

import (
	"database/sql"

	"github.com/andrewcharlton/school-dashboard/analysis/subject"
)

// Connect to the database
func Connect(filename string) (Database, error) {

	var db Database

	// Create connection to the database
	conn, err := sql.Open("sqlite3", filename)
	if err != nil {
		return Database{}, err
	}
	db.conn = conn

	// Load cached variables
	err = db.Refresh()
	if err != nil {
		return Database{}, err
	}

	return db, nil
}

// Close the connection to the database.
func (db Database) Close() error {

	err := db.conn.Close()
	return err
}

// Refresh the cached config from the database
func (db *Database) Refresh() error {

	return nil
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
func (db *Database) loadGrades(lvl int) (map[string]subject.Grade, error) {

	rows, err := db.conn.Query(`SELECT grade, points, att8, l1_pass, l2_pass
								FROM grades
								WHERE level_id=?`, lvl)
	if err != nil {
		return map[string]subject.Grade{}, err
	}
	defer rows.Close()

	grades := map[string]subject.Grade{}
	for rows.Next() {
		var g subject.Grade
		err := rows.Scan(&g.Grd, &g.Pts, &g.Att8, &g.L1Pass, &g.L2Pass)
		if err != nil {
			return map[string]subject.Grade{}, err
		}
		grades[g.Grd] = g
	}
	return grades, nil
}

// loadTMs loads up the transition matrices for a particular subject
func (db Database) loadTMs(subj *subject.Subject, tmName string) error {

	rows, err := db.conn.Query(`SELECT year, ks2, grade, probability FROM tms
								WHERE subject=?`, tmName)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var year, ks2, grade string
		var prob float64
		err := rows.Scan(&year, &ks2, &grade, &prob)
		if err != nil {
			return err
		}
		tm, exists := subj.TMs[year]
		if !exists {
			tm = subject.NewTM(subj.Level)
		}
		row, exists := tm.Rows[ks2]
		if !exists {
			row = map[string]float64{}
		}
		row[grade] = prob

		tm.Rows[ks2] = row
		subj.TMs[year] = tm
	}

	return nil
}

// loadSubjects pulls in the subject list, indexed by subject id
func (db *Database) loadSubjects() error {

	rows, err := db.conn.Query(`SELECT id, subject, level_id, ebacc, ks2_prior, tm
								FROM subjects`)
	if err != nil {
		return err
	}
	defer rows.Close()

	db.subjects = map[int]*subject.Subject{}
	for rows.Next() {
		var lvl int
		var tm string
		var s subject.Subject
		err := rows.Scan(&s.SubjID, &s.Subj, &lvl, &s.EBacc, &s.KS2Prior, &tm)
		if err != nil {
			return err
		}
		s.Level = db.levels[lvl]

		err = db.loadTMs(&s, tm)
		if err != nil {
			return err
		}

		db.subjects[s.SubjID] = &s
	}
	return nil
}
