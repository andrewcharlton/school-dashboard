package database

import (
	"database/sql"

	"github.com/andrewcharlton/school-dashboard/analysis/student"
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

	err := db.loadConfig()
	if err != nil {
		return err
	}

	err = db.loadDates()
	if err != nil {
		return err
	}

	err = db.loadResultsets()
	if err != nil {
		return err
	}

	err = db.loadNatYears()
	if err != nil {
		return err
	}

	err = db.loadEthnicities()
	if err != nil {
		return err
	}

	err = db.loadLevels()
	if err != nil {
		return err
	}

	err = db.loadSubjects()
	if err != nil {
		return err
	}

	err = db.loadAttainment8()
	if err != nil {
		return err
	}

	err = db.prepareStatements()
	if err != nil {
		return err
	}

	return nil
}

// Dates returns a sorted list of all Dates in the database that
// are marked to be listed.
func (db *Database) loadDates() error {

	rows, err := db.conn.Query(`SELECT id, date FROM dates WHERE list=1
								ORDER BY id`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var l Lookup
		err := rows.Scan(&l.ID, &l.Name)
		if err != nil {
			return err
		}
		db.Dates = append(db.Dates, l)
		db.dateMap[l.ID] = l.Name
	}

	return nil
}

// Ethnicities returns all the distinct ethnicities present in
// the database, in descending order of number of students present
func (db *Database) loadEthnicities() error {

	rows, err := db.conn.Query(`SELECT ethnicity, COUNT(1) as n
								FROM students
								GROUP BY ethnicity
								ORDER BY n DESC`)
	if err != nil {
		return err
	}
	defer rows.Close()

	eth := []string{}
	for rows.Next() {
		var name string
		var count int
		err := rows.Scan(&name, &count)
		if err != nil {
			return err
		}
		eth = append(eth, name)
	}

	for n, e := range eth {
		switch {
		case n < 8:
			db.Ethnicities = append(db.Ethnicities, e)
		default:
			db.OtherEths[e] = true
		}
	}
	db.Ethnicities = append(db.Ethnicities, "Other")

	return nil
}

// grades pulls in the list of grades at a particular level
func (db Database) loadGrades(lvl int) (map[string]subject.Grade, error) {

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

// levels pulls in all of the levels for caching
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

// natYears returns a list of the national years
func (db *Database) loadNatYears() error {

	rows, err := db.conn.Query(`SELECT id, year FROM nat_years`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var l Lookup
		err := rows.Scan(&l.ID, &l.Name)
		if err != nil {
			return err
		}
		db.NatYears = append(db.NatYears, l)
		db.natYearMap[l.ID] = l.Name
	}

	return nil
}

// progress8 scores for each year
func (db *Database) loadAttainment8() error {

	rows, err := db.conn.Query(`SELECT year, ks2, att8, english, maths, ebacc, other
								FROM nat_progress8
								INNER JOIN ON year_id = nat_years.id`)
	if err != nil {
		return err
	}
	defer rows.Close()

	db.Attainment8 = map[string](map[string]student.Attainment8){}
	for rows.Next() {
		var year, ks2 string
		var a8, en, ma, eb, oth float64
		err = rows.Scan(&year, &ks2, &a8, &en, &ma, &eb, &oth)
		if err != nil {
			return err
		}

		att8, exists := db.Attainment8[year]
		if !exists {
			att8 = map[string]student.Attainment8{}
		}
		att8[ks2] = student.Attainment8{English: en, Maths: ma, EBacc: eb,
			Other: oth, Overall: a8}
		db.Attainment8[year] = att8
	}

	return nil
}

// Resultsets returns a sorted list of all Resultsets in the database.
// An 'Exams Only' option encapsulates all individual exam resultsets,
// all other resultsets marked to be listed are included.
func (db *Database) loadResultsets() error {

	rows, err := db.conn.Query(`SELECT id, resultset FROM resultsets
								WHERE is_exam=0 AND list=1
								ORDER BY id`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var l Lookup
		err := rows.Scan(&l.ID, &l.Name)
		if err != nil {
			return err
		}
		db.Resultsets = append(db.Resultsets, l)
		db.resultsetMap[l.ID] = l.Name
	}

	return nil
}

// subjects pulls in the subject list, indexed by subject id
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

		err = db.tms(&s, tm)
		if err != nil {
			return err
		}
		db.subjects[s.SubjID] = &s
	}
	return nil
}

// tms loads up the transition matrices for a particular subject
func (db Database) tms(subj *subject.Subject, tmName string) error {

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

// prepareStatements prepares all query statements ready for use
func (db *Database) prepareStatements() error {

	// Close any existing statements
	for _, s := range db.stmts {
		err := s.Close()
		if err != nil {
			return err
		}
	}

	// Prepare statements
	db.stmts = map[string]*sql.Stmt{}
	for key, sql := range sqlStatements {
		s, err := db.conn.Prepare(sql)
		if err != nil {
			return err
		}
		db.stmts[key] = s
	}
	return nil
}
