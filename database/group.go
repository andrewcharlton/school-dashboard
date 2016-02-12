package database

import (
	"fmt"

	"github.com/andrewcharlton/school-dashboard/analysis/group"
	"github.com/andrewcharlton/school-dashboard/analysis/student"
)

// GroupByFilter returns a group of students who satisfy the criteria
// specified in the filter.
func (db Database) GroupByFilter(f Filter) (group.Group, error) {

	query := f.sql("students")
	upns, err := db.getUPNs(query)
	if err != nil {
		return group.Group{}, err
	}

	return db.group(upns, f)
}

// GroupByClass returns a group of students who are present in the
// subject/class specified at the particular date.
func (db Database) GroupByClass(subjID, class string, f Filter) (group.Group, error) {

	rows, err := db.stmts["inClass"].Query(f.Date, subjID, class)
	if err != nil {
		return group.Group{}, err
	}
	defer rows.Close()

	upns := []string{}
	for rows.Next() {
		var upn string
		if err := rows.Scan(&upn); err != nil {
			return group.Group{}, err
		}
		upns = append(upns, upn)
	}

	return db.group(upns, f)
}

// GroupByFilteredClass returns a group of students who meet the filter criteria
// and are in the subject/class combination specified.
// If class="", then all students in that subject are returned.
func (db Database) GroupByFilteredClass(subjID, class string, f Filter) (group.Group, error) {

	if subjID == "" {
		return db.GroupByFilter(f)
	}

	query := f.sql("classes_filter")
	query += fmt.Sprintf(` AND subject_id=%v`, subjID)
	if class != "" {
		query += fmt.Sprintf(` AND class="%v"`, class)
	}

	upns, err := db.getUPNs(query)
	if err != nil {
		return group.Group{}, err
	}

	return db.group(upns, f)
}

// group returns a list of students loaded based upon the supplied
// filters.
func (db Database) group(upns []string, f Filter) (group.Group, error) {

	students := []student.Student{}

	for _, upn := range upns {
		student, err := db.Student(upn, f)
		if err != nil {
			return group.Group{}, err
		}
		students = append(students, student)
	}

	return group.Group{students}, nil
}

// getUPNs queries the database and returns a list of UPNs from
func (db Database) getUPNs(query string) ([]string, error) {

	// Ensure group is sorted
	query += ` ORDER BY (surname || " " || forename)`

	rows, err := db.conn.Query(query)
	if err != nil {
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
