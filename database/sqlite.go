package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/andrewcharlton/school-dashboard/analysis/national"
	_ "github.com/mattn/go-sqlite3" // SQL Driver
)

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

// Level returns the named level
func (db sqliteDB) Level(name string) *lvl.Level {

	for _, l := range db.levels {
		if l.Lvl == name {
			return l
		}
	}
	return nil
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

// GroupByFilter returns a group of students who satisfy the criteria
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
						  WHERE date_id=%v`,
		table, f.Date)

	if f.Year != "" {
		query += fmt.Sprintf(" AND year=%v", f.Year)
	}
	if f.PP != "" {
		query += fmt.Sprintf(" AND pp=%v", f.PP)
	}
	if f.EAL != "" {
		query += fmt.Sprintf(" AND eal=%v", f.EAL)
	}
	if f.Gender != "" {
		query += fmt.Sprintf(" AND gender=%v", f.Gender)
	}
	if len(f.SEN) > 0 {
		query += fmt.Sprintf(` AND sen_status IN ("` + strings.Join(f.SEN, `", "`) + `")`)
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

	prog8, err := db.loadProgress8(yearID)
	if err != nil {
		return national.National{}, err
	}

	tms, err := db.loadTMs(yearID)
	if err != nil {
		return national.National{}, err
	}

	nat := national.National{Prog8: prog8, TMs: tms}
	return nat, nil
}

func (db sqliteDB) loadProgress8(yearID string) (map[string]national.Progress8, error) {

	// Load attainment 8 data
	rows, err := db.conn.Query(`SELECT ks2, att8, english, maths, ebacc, other
								FROM nat_progress8
								WHERE year_ID=?`, yearID)
	if err != nil {
		return map[string]national.Progress8{}, err
	}
	defer rows.Close()

	prog8 := map[string]national.Progress8{}
	for rows.Next() {
		var ks2 string
		var a8, en, ma, eb, oth float64
		err := rows.Scan(&ks2, &a8, &en, &ma, &eb, &oth)
		if err != nil {
			return map[string]national.Progress8{}, err
		}
		prog8[ks2] = national.Progress8{English: en, Maths: ma, EBacc: eb,
			Other: oth, Att8: a8}
	}

	return prog8, nil
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
