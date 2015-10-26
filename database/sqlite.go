package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/andrewcharlton/school-dashboard/analysis"
	_ "github.com/mattn/go-sqlite3" // SQL Driver
)

// A sqliteDB is a container for the
type sqliteDB struct {
	conn  *sql.DB              // Database connection
	stmts map[string]*sql.Stmt // Prepared statements for quicker queries

	// Cached objects
	subjects map[int]*analysis.Subject
	levels   map[int]*analysis.Level
}

// ConnectSQLite opens a connection to the database
// and returns a Database object
func ConnectSQLite(filename string) (Database, error) {

	var DB Database

	conn, err := sql.Open("sqlite3", filename)
	if err != nil {
		return DB, err
	}

	var db sqliteDB
	db.conn = conn

	// Prepare SQL statements
	if err := db.prepareStatements(); err != nil {
		return DB, err
	}

	// Load details of levels cache
	if err := db.loadLevels(); err != nil {
		return DB, err
	}

	// Load subjects to cache
	if err := db.loadSubjects(); err != nil {
		return DB, err
	}

	DB = db
	return DB, nil
}

// Close terminates the connection to the
// database.
func (db sqliteDB) Close() error {

	err := db.conn.Close()
	return err
}

// Query strings for prepared statements
var sqliteStmts = map[string]string{

	"student": `SELECT upn, surname, forename, year, form,
				pp, eal, gender, ethnicity, sen_status,
				sen_info, sen_strat, ks2_aps, ks2_band,
				ks2_en, ks2_ma, ks2_av
				FROM students
				WHERE upn=? AND date=?`,

	"courses": `SELECT subject_id, subject, grade FROM results
				WHERE upn=? AND resultset=?`,
}

// prepareStatements prepares a query statement for each sql string
// storedin sqliteStmnts
func (db *sqliteDB) prepareStatements() error {

	// Close any existing statements
	for _, stmt := range db.stmts {
		if err := stmt.Close(); err != nil {
			return err
		}
	}

	// Recreate statements
	db.stmts = map[string]*sql.Stmt{}
	for key, sql := range sqliteStmts {
		stmt, err := db.conn.Prepare(sql)
		if err != nil {
			return err
		}
		db.stmts[key] = stmt
	}
	return nil
}

// loadLevels pulls in all of the levels for caching
func (db *sqliteDB) loadLevels() error {

	rows, err := db.conn.Query("SELECT id, level, is_gcse FROM levels")
	if err != nil {
		return err
	}
	defer rows.Close()

	db.levels = map[int]*analysis.Level{}
	for rows.Next() {
		var id int
		var l analysis.Level
		err := rows.Scan(&id, &l.Lvl, &l.IsGCSE)
		if err != nil {
			return err
		}

		grades, err := db.loadGrades(id)
		if err != nil {
			return err
		}
		l.Gradeset = grades

		db.levels[id] = &l
	}

	return nil
}

// loadGrades pulls in the list of grades at a particular level
func (db *sqliteDB) loadGrades(level int) (map[string]*analysis.Grade, error) {

	rows, err := db.conn.Query(`SELECT grade, points, att8, l1pass, l2pass
								FROM grades
								WHERE level_id=?`, level)
	if err != nil {
		return map[string]*analysis.Grade{}, err
	}
	defer rows.Close()

	grades := map[string]*analysis.Grade{}
	for rows.Next() {
		var g analysis.Grade
		err := rows.Scan(&g.Grd, &g.Pts, &g.Att8, &g.L1Pass, &g.L2Pass)
		if err != nil {
			return map[string]*analysis.Grade{}, err
		}
		grades[g.Grd] = &g
	}

	return grades, nil
}

// loadSubjects pulls in the subject list for caching
func (db *sqliteDB) loadSubjects() error {

	rows, err := db.conn.Query(`SELECT id, subject, level_id, ebacc, ks2_prior
								FROM subjects`)
	if err != nil {
		return err
	}
	defer rows.Close()

	db.subjects = map[int]*analysis.Subject{}
	for rows.Next() {
		var id, level int
		var s analysis.Subject
		err := rows.Scan(&id, &s.Subj, &level, &s.EBacc, &s.KS2Prior)
		if err != nil {
			return err
		}
		s.Level = db.levels[level]
		db.subjects[id] = &s
	}

	return nil
}

// Dates returns a sorted list of all Dates in the database that
// are marked to be listed.
func (db sqliteDB) Dates() ([]Lookup, error) {

	rows, err := db.conn.Query(`SELECT id, date FROM dates WHERE list=1
								ORDER BY id`)
	if err != nil {
		return []Lookup{}, err
	}
	defer rows.Close()

	dates := []Lookup{}
	for rows.Next() {
		var l Lookup
		err := rows.Scan(&l.ID, &l.Name)
		if err != nil {
			return []Lookup{}, err
		}
		dates = append(dates, l)
	}

	return dates, nil
}

// Resultsets returns a sorted list of all Resultsets in the database.
// An 'Exams Only' option encapsulates all individual exam resultsets,
// all other resultsets marked to be listed are included.
func (db sqliteDB) Resultsets() ([]Lookup, error) {

	rows, err := db.conn.Query(`SELECT id, resultset FROM resultsets
								WHERE is_exam=0 AND list=1
								ORDER BY id`)
	if err != nil {
		return []Lookup{}, err
	}
	defer rows.Close()

	rs := []Lookup{{"0", "Exams Only"}}
	for rows.Next() {
		var l Lookup
		err := rows.Scan(&l.ID, &l.Name)
		if err != nil {
			return []Lookup{}, err
		}
		rs = append(rs, l)
	}

	return rs, nil
}

// Ethnicities returns all the distinct ethnicities present in
// the database, and the frequency that each appears with.
// Only students who are present in listed dates are counted.
func (db sqliteDB) Ethnicities() ([]Ethnicity, error) {

	rows, err := db.conn.Query(`SELECT ethnicity, COUNT(1) as n
								FROM students
								GROUP BY ethnicity`)
	if err != nil {
		return []Ethnicity{}, err
	}
	defer rows.Close()

	eth := []Ethnicity{}
	for rows.Next() {
		var e Ethnicity
		err := rows.Scan(&e.Name, &e.Count)
		if err != nil {
			return []Ethnicity{}, err
		}
		eth = append(eth, e)
	}

	return eth, nil
}

// Group returns a list of UPNs for students who satisfy the criteria
// specified in the filter.
func (db sqliteDB) Group(f Filter) (analysis.Group, error) {

	query := fmt.Sprintf(`SELECT upn FROM students WHERE date = %v`, f.Date)

	if len(f.Years) > 0 {
		query += " AND (year IN (" + strings.Join(f.Years, ", ") + ")"
	}
	if f.PP != "" {
		query += fmt.Sprintf(" AND pp = %v", f.PP)
	}
	if f.EAL != "" {
		query += fmt.Sprintf(" AND eal = %v", f.EAL)
	}
	if f.Gender != "" {
		query += fmt.Sprintf(" AND gender = %v", f.Gender)
	}
	if len(f.SEN) > 0 {
		query += fmt.Sprintf(" AND sen IN (" + strings.Join(f.SEN, ", ") + ")")
	}
	if len(f.Ethnicities) > 0 {
		query += fmt.Sprintf(" AND ethnicity IN (" + strings.Join(f.Ethnicities, ", ") + ")")
	}
	if len(f.KS2Bands) > 0 {
		query += fmt.Sprintf(" AND ks2_band IN (" + strings.Join(f.KS2Bands, ", ") + ")")
	}
	rows, err := db.conn.Query(query)
	if err != nil {
		return analysis.Group{}, err
	}
	defer rows.Close()

	upns := []string{}
	for rows.Next() {
		var upn string
		if err := rows.Scan(&upn); err != nil {
			return analysis.Group{}, err
		}
		upns = append(upns, upn)
	}

	students := []analysis.Student{}
	for _, upn := range upns {
		sf := StudentFilter{upn, f.Date, f.Resultset}
		student, err := db.Student(sf)
		if err != nil {
			return analysis.Group{}, err
		}
		students = append(students, student)
	}

	return analysis.Group{students}, nil
}

// Student creates a student object based on the
func (db sqliteDB) Student(f StudentFilter) (analysis.Student, error) {

	row := db.stmts["student"].QueryRow(f.UPN, f.Date)

	sen := analysis.SENInfo{}
	ks2 := analysis.KS2Info{}
	s := analysis.Student{}
	err := row.Scan(&s.UPN, &s.Surname, &s.Forename, &s.Year, &s.Form,
		&s.PP, &s.EAL, &s.Gender, &s.Ethnicity, &sen.Status, &sen.Info,
		&sen.Strategies, &ks2.APS, &ks2.Band, &ks2.En,
		&ks2.Ma, &ks2.Av)
	if err != nil {
		return analysis.Student{}, err
	}
	s.SEN = sen
	s.KS2 = ks2

	// Load courses
	rows, err := db.stmts["courses"].Query(f.UPN, f.Resultset)
	if err != nil {
		return analysis.Student{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var subjID int
		var subj, grade string
		err := rows.Scan(&subjID, &subj, &grade)
		if err != nil {
			return analysis.Student{}, err
		}

		subject := db.subjects[subjID]
		c := analysis.Course{subject, subject.Gradeset[grade]}
		s.Courses[subj] = c
	}

	return s, nil
}
