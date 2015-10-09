package database

import (
	"database/sql"
	"fmt"
	"strings"
)

// Filter holds t
type Filter struct {
	// Effective date to use for membership of the group
	Date string
	// Which set of assessment results to use
	Resultset string
	// How to cope with exams: "First", "Best", "Ignore"
	Exams string
	// Which yeargroups to include in the group
	Years []string
	// Pupil premium: "", "1", or "0" for Any/True/False
	PP string
	// EAL students: "", "1", or "0" for Any/True/False
	EAL string
	// Gender: "", "1", "0" for Any/Male/Female
	Gender string
	// SEN types to include - Empty for any
	SEN []string
	// Which Ethnic groups to include, empty for any
	Ethnicities []string
	// Which KS2 bands to include, empty for any
	KS2Bands []string
}

// StudentSQL generates the SQL string needed to query the
// database.  Querying the database with this string produces
// a list of UPN's of students who meet the criteria set out
// in the filter.
func (f Filter) StudentSQL() string {

	query := fmt.Sprintf(`SELECT upn FROM students WHERE date = %v`, f.Date)

	if len(f.Years) > 0 {
		query += " AND (year = "
		query += strings.Join(f.Years, " OR year = ")
		query += ")"
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

	return query + " ORDER BY students.surname, students.forename"
}

// DefaultFilter selects the most recent effective date and resultset
// and creates Filter object with those.
func (db SchoolDB) DefaultFilter() (Filter, error) {

	f := Filter{}

	err := db.QueryRow("SELECT id FROM dates ORDER BY date DESC LIMIT 1").Scan(&f.Date)
	if err != nil {
		return f, err
	}

	err = db.QueryRow("SELECT id FROM resultsets WHERE is_exam=0 ORDER BY date DESC LIMIT 1").Scan(&f.Resultset)
	switch {
	case err == sql.ErrNoRows:
		f.Resultset = "0"
	case err != nil:
		return f, err
	}

	// Select all Ethnicities
	rows, err := db.Query("SELECT DISTINCT ethnicity FROM students")
	if err != nil {
		return f, err
	}
	defer rows.Close()
	for rows.Next() {
		var eth string
		err := rows.Scan(&eth)
		if err != nil {
			return f, err
		}
		f.Ethnicities = append(f.Ethnicities, eth)
	}
	f.Ethnicities = append(f.Ethnicities, "Other")

	f.Exams = "First"
	f.Years = []string{"11"}
	f.KS2Bands = []string{"High", "Middle", "Low", "None"}
	f.SEN = []string{"", "S", "A", "P"}
	return f, nil
}
