package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/andrewcharlton/school-dashboard/analysis"
	"github.com/andrewcharlton/school-dashboard/level"
	"github.com/andrewcharlton/school-dashboard/national"
	_ "github.com/mattn/go-sqlite3" // SQL Driver
)

// A sqliteDB is a container for the
type sqliteDB struct {
	conn  *sql.DB              // Database connection
	stmts map[string]*sql.Stmt // Prepared statements for quicker queries

	// Cached objects
	subjects map[int]*analysis.Subject
	levels   map[int]*level.Level
	years    map[string]int // Map of dates to school_years
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

	// Load year data for attendance
	if err := db.loadYears(); err != nil {
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
				sen_need, sen_info, sen_strat, sen_access, sen_iep, 
				ks2_aps, ks2_band, ks2_en, ks2_ma, ks2_av, ks2_re, 
				ks2_wr, ks2_gps
				FROM students
				WHERE upn=? AND date_id=?`,

	"results": `SELECT subject_id, subject, grade, effort FROM results
				WHERE upn=? AND resultset=?`,

	"classes": `SELECT subject_id, class, teacher FROM classes
				WHERE upn=? AND date_id=?
				ORDER BY class`,

	"inClass": `SELECT upn FROM classes
				WHERE date_id=? AND subject_id=? AND class=?
				ORDER BY name`,

	"subjects": `SELECT subjects.id as id
				FROM subjects
				INNER JOIN levels ON subjects.level_id = levels.id
				WHERE keystage=?`,

	"classlist": `SELECT DISTINCT class FROM classes
					WHERE date_id=? AND subject_id=?
					ORDER BY (
						CASE WHEN substr(class, 1, 1) == "1"
						THEN class
						ELSE  "0" || class
					END)`,

	"bestExams": `SELECT subject_id, subject, grade FROM results
				  WHERE upn=? AND is_exam=1
				  ORDER BY points DESC`,

	"firstExams": `SELECT subject_id, subject, grade FROM results
				  WHERE upn=? AND is_exam=1
				  ORDER BY date`,

	"search": `SELECT upn, surname, forename, year, form
			   FROM students
			   WHERE date_id=? AND
			   (((forename || " " || surname) LIKE ?)
				OR ((surname || ", " || forename) LIKE ?))
				ORDER BY (surname || " " || forename);`,

	"attendance": `SELECT poss_year, absence_year, unauth_year, mon_am,
					mon_pm, tue_am, tue_pm, wed_am, wed_pm, thu_am,
					thu_pm, fri_am, fri_pm
					FROM attendance
					WHERE upn=? AND year_id=?
					ORDER BY week_start DESC
					LIMIT 1`,
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

	db.levels = map[int]*level.Level{}
	for rows.Next() {
		var id int
		var l level.Level
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
func (db *sqliteDB) loadGrades(lvl int) (map[string]*level.Grade, error) {

	rows, err := db.conn.Query(`SELECT grade, points, att8, l1_pass, l2_pass
								FROM grades
								WHERE level_id=?`, lvl)
	if err != nil {
		return map[string]*level.Grade{}, err
	}
	defer rows.Close()

	grades := map[string]*level.Grade{}
	for rows.Next() {
		var g level.Grade
		err := rows.Scan(&g.Grd, &g.Pts, &g.Att8, &g.L1Pass, &g.L2Pass)
		if err != nil {
			return map[string]*level.Grade{}, err
		}
		grades[g.Grd] = &g
	}

	return grades, nil
}

// loadSubjects pulls in the subject list for caching
func (db *sqliteDB) loadSubjects() error {

	rows, err := db.conn.Query(`SELECT id, subject, level_id, ebacc, tm, ks2_prior
								FROM subjects`)
	if err != nil {
		return err
	}
	defer rows.Close()

	db.subjects = map[int]*analysis.Subject{}
	for rows.Next() {
		var id, lvl int
		var s analysis.Subject
		err := rows.Scan(&id, &s.Subj, &lvl, &s.EBacc, &s.TM, &s.KS2Prior)
		if err != nil {
			return err
		}
		s.Level = db.levels[lvl]
		db.subjects[id] = &s
	}

	return nil
}

// loadYears pulls in the effective dates and the corresponding year
func (db *sqliteDB) loadYears() error {

	rows, err := db.conn.Query(`SELECT id, year_id FROM dates`)
	if err != nil {
		return err
	}
	defer rows.Close()

	db.years = map[string]int{}
	for rows.Next() {
		var date string
		var year int
		err := rows.Scan(&date, &year)
		if err != nil {
			return err
		}
		db.years[date] = year
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

// Config retrieves the config values from the database
func (db sqliteDB) Config() (Config, error) {

	rows, err := db.conn.Query("SELECT key, value FROM config")
	if err != nil {
		return Config{}, err
	}
	defer rows.Close()

	cfg := Config{}
	for rows.Next() {
		var key, value string
		err := rows.Scan(&key, &value)
		if err != nil {
			return Config{}, err
		}
		switch key {
		case "School":
			cfg.School = value
		case "URN":
			cfg.URN = value
		case "LEA":
			cfg.LEA = value
		case "Date":
			cfg.Date = value
		case "Resultset":
			cfg.Resultset = value
		case "NatYear":
			cfg.NatYear = value
		case "Year":
			cfg.Year = value
		}
	}
	return cfg, nil
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

	rs := []Lookup{}
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

// Subjects returns all of the subjects, indexed by their id number
func (db sqliteDB) Subjects() map[int]*analysis.Subject {
	return db.subjects
}

// Classes returns a sorted list of classes present
// for a subject.
func (db sqliteDB) Classes(subj_id, date string) ([]string, error) {

	rows, err := db.stmts["classlist"].Query(date, subj_id)
	if err == sql.ErrNoRows {
		return []string{}, errors.New("No classes present for this subject, on this date.")
	}
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	classes := []string{}
	for rows.Next() {
		var class string
		if err := rows.Scan(&class); err != nil {
			return []string{}, err
		}
		classes = append(classes, class)
	}

	return classes, nil
}

// Ethnicities returns all the distinct ethnicities present in
// the database, and the frequency that each appears with.
// Only students who are present in listed dates are counted.
func (db sqliteDB) Ethnicities() ([]Ethnicity, error) {

	rows, err := db.conn.Query(`SELECT ethnicity, COUNT(1) as n
								FROM students
								GROUP BY ethnicity
								ORDER BY n DESC`)
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

// group returns a list of students loaded based upon the supplied
// filters.
func (db sqliteDB) group(upns []string, f StudentFilter) (analysis.Group, error) {

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

// Group returns a group of students who satisfy the criteria
// specified in the filter.
func (db sqliteDB) GroupByFilter(f Filter) (analysis.Group, error) {

	query := db.groupFilter(f, "students")
	query += ` ORDER BY (surname || " " || forename)`

	upns, err := db.getUPNs(query)
	if err != nil {
		return analysis.Group{}, err
	}

	return db.group(upns, StudentFilter{"", f.Date, f.Resultset})
}

// groupFilter dynamically constructs the SQL statement from a filter.
func (db sqliteDB) groupFilter(f Filter, table string) string {

	query := fmt.Sprintf(`SELECT upn FROM %v
						  WHERE date_id = %v and year= %v`,
		table, f.Date, f.Year)

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
		query += fmt.Sprintf(` AND sen IN ("` + strings.Join(f.SEN, `", "`) + `")`)
	}
	if len(f.Ethnicities) > 0 {
		query += fmt.Sprintf(` AND ethnicity IN ("` + strings.Join(f.Ethnicities, `", "`) + `")`)
	}
	if len(f.KS2Bands) > 0 {
		query += fmt.Sprintf(` AND ks2_band IN ("` + strings.Join(f.KS2Bands, `", "`) + `")`)
	}

	return query
}

func (db sqliteDB) getUPNs(query string) ([]string, error) {

	rows, err := db.conn.Query(query)
	if err != nil {
		fmt.Println(query)
		return []string{}, err
	}
	defer rows.Close()

	upns := []string{}
	for rows.Next() {
		var upn string
		if err := rows.Scan(&upn); err != nil {
			return []string{}, err
		}
		upns = append(upns, upn)
	}
	return upns, nil
}

// GroupByClass returns a group of students who are present in the
// subject/class specified at the particular date.
func (db sqliteDB) GroupByClass(subj_id, class string, f Filter) (analysis.Group, error) {

	rows, err := db.stmts["inClass"].Query(f.Date, subj_id, class)
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

	return db.group(upns, StudentFilter{"", f.Date, f.Resultset})
}

// GroupByFilteredClass returns a group of students who meet the filter criteria
// and are in the subject/class combination specified.
// If class="", then all students in that subject are returned.
func (db sqliteDB) GroupByFilteredClass(subj_id, class string, f Filter) (analysis.Group, error) {

	if subj_id == "" {
		return db.GroupByFilter(f)
	}

	query := db.groupFilter(f, "classes_filter")
	query += fmt.Sprintf(` AND subject_id=%v`, subj_id)
	if class != "" {
		query += fmt.Sprintf(` AND class="%v"`, class)
	}
	query += ` ORDER BY (surname || " " || forename)`

	upns, err := db.getUPNs(query)
	if err != nil {
		return analysis.Group{}, err
	}

	return db.group(upns, StudentFilter{"", f.Date, f.Resultset})
}

// Student creates a student object based on the
func (db sqliteDB) Student(f StudentFilter) (analysis.Student, error) {

	// Load student details
	row := db.stmts["student"].QueryRow(f.UPN, f.Date)
	sen := analysis.SENInfo{}
	ks2 := analysis.KS2Info{}
	var male bool
	s := analysis.Student{}
	err := row.Scan(&s.UPN, &s.Surname, &s.Forename, &s.Year, &s.Form,
		&s.PP, &s.EAL, &male, &s.Ethnicity, &sen.Status, &sen.Need,
		&sen.Info, &sen.Strategies, &sen.Access, &sen.IEP, &ks2.APS,
		&ks2.Band, &ks2.En, &ks2.Ma, &ks2.Av, &ks2.Re, &ks2.Wr, &ks2.GPS)
	if err == sql.ErrNoRows {
		return analysis.Student{}, errors.New("Student not on roll at this date.  Try changing the date, or search for another student.")
	}
	if err != nil {
		return analysis.Student{}, err
	}
	s.SEN = sen
	s.KS2 = ks2
	if male {
		s.Gender = "Male"
	} else {
		s.Gender = "Female"
	}

	// Load courses
	var rows *sql.Rows
	switch f.Resultset {
	case "1":
		rows, err = db.stmts["firstExams"].Query(f.UPN)
	case "2":
		rows, err = db.stmts["bestExams"].Query(f.UPN)
	default:
		rows, err = db.stmts["results"].Query(f.UPN, f.Resultset)
	}

	if err != nil {
		return analysis.Student{}, err
	}
	defer rows.Close()

	s.Courses = map[string]analysis.Course{}
	for rows.Next() {
		var subjID, effort int
		var subj, grade string
		err := rows.Scan(&subjID, &subj, &grade, &effort)
		if err != nil {
			return analysis.Student{}, err
		}

		// Don't overwrite existing course
		_, exists := s.Courses[subj]
		if exists {
			continue
		}

		subject := db.subjects[subjID]
		c := analysis.Course{subject, subject.Gradeset[grade], effort, "", ""}
		s.Courses[subj] = c
	}

	// Load class data
	rows, err = db.stmts["classes"].Query(f.UPN, f.Date)
	if err != nil {
		return analysis.Student{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var subj_id int
		var class, teacher string
		err := rows.Scan(&subj_id, &class, &teacher)
		if err != nil {
			return analysis.Student{}, nil
		}
		subj := db.subjects[subj_id].Subj
		if c, exists := s.Courses[subj]; exists {
			c.Class = class
			c.Teacher = teacher
			s.Courses[subj] = c
		}
	}

	// Load attendance data
	row = db.stmts["attendance"].QueryRow(f.UPN, f.Date)
	att := analysis.AttendanceInfo{}
	err = row.Scan(&att.Possible, &att.Absences, &att.Unauthorised,
		&att.Sessions[0], &att.Sessions[1], &att.Sessions[2],
		&att.Sessions[3], &att.Sessions[4], &att.Sessions[5],
		&att.Sessions[6], &att.Sessions[7], &att.Sessions[8],
		&att.Sessions[9])
	if err != nil && err != sql.ErrNoRows {
		return analysis.Student{}, nil
	}
	s.Attendance = att

	return s, nil
}

// Search returns a list of students that satisfy the search criteria.
func (db sqliteDB) Search(name, date string) ([]StudentLookup, error) {

	str := "%" + strings.Replace(name, "*", "%", -1) + "%"
	rows, err := db.stmts["search"].Query(date, str, str)
	if err != nil {
		return []StudentLookup{}, err
	}
	defer rows.Close()

	list := []StudentLookup{}
	for rows.Next() {
		var upn, surname, forename, year, form string
		err := rows.Scan(&upn, &surname, &forename, &year, &form)
		if err != nil {
			return []StudentLookup{}, err
		}
		s := StudentLookup{UPN: upn,
			Name: surname + ", " + forename,
			Year: year,
			Form: form}
		list = append(list, s)
	}

	return list, nil
}

// NationalYears returns a list of years where national data is
// available in the database.
func (db sqliteDB) NationalYears() ([]Lookup, error) {

	rows, err := db.conn.Query(`SELECT id, year FROM nat_years`)
	if err != nil {
		return []Lookup{}, err
	}
	defer rows.Close()

	years := []Lookup{}
	for rows.Next() {
		var id, year string
		err := rows.Scan(&id, &year)
		if err != nil {
			return []Lookup{}, err
		}
		years = append(years, Lookup{id, year})
	}

	return years, nil
}

// National returns a set of national data for the given year.
func (db sqliteDB) National(yearID string) (national.National, error) {

	att8, err := db.loadAtt8(yearID)
	if err != nil {
		return national.National{}, err
	}

	tms, err := db.loadTMs(yearID)
	if err != nil {
		return national.National{}, err
	}

	nat := national.National{Att8: att8, TMs: tms}
	return nat, nil
}

func (db sqliteDB) loadAtt8(yearID string) (map[string]float64, error) {

	// Load attainment 8 data
	rows, err := db.conn.Query(`SELECT ks2, att8 FROM nat_progress8
								WHERE year_ID=?`, yearID)
	if err != nil {
		return map[string]float64{}, err
	}
	defer rows.Close()

	att8 := map[string]float64{}
	for rows.Next() {
		var ks2 string
		var a8 float64
		err := rows.Scan(&ks2, &a8)
		if err != nil {
			return map[string]float64{}, err
		}
		att8[ks2] = a8
	}

	return att8, nil
}

// loadTMs loads up the transition matrices from a particular year
func (db sqliteDB) loadTMs(yearID string) (map[string]national.TransitionMatrix, error) {

	// Load Transition Matrices
	rows, err := db.conn.Query(`SELECT subject, level_id, ks2, grade, probability
								FROM tms
								WHERE year_ID=?`, yearID)
	if err != nil {
		return map[string]national.TransitionMatrix{}, err
	}
	defer rows.Close()

	tms := map[string]national.TransitionMatrix{}
	for rows.Next() {
		var subj, ks2, grade string
		var lvl int
		var prob float64
		err := rows.Scan(&subj, &lvl, &ks2, &grade, &prob)
		if err != nil {
			fmt.Println("loadTMs Error: ", err)
			return map[string]national.TransitionMatrix{}, err
		}
		tm, exists := tms[subj]
		if !exists {
			tmRows := map[string]national.TMRow{}
			tm = national.TransitionMatrix{Rows: tmRows, Level: db.levels[lvl]}
		}
		row, exists := tm.Rows[ks2]
		if !exists {
			row = national.TMRow{}
		}
		row[grade] = prob
		tm.Rows[ks2] = row
		tms[subj] = tm
	}

	return tms, nil
}
