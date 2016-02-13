package database

import (
	"fmt"
	"strings"
)

// CurrentWeek looks up the name of the latest week, as used
// for attendance.
func (db Database) CurrentWeek() (string, error) {

	var week string
	row := db.stmts["currentWeek"].QueryRow()
	err := row.Scan(&week)
	if err != nil {
		return "", err
	}

	digits := strings.Split(week, "-")
	if len(digits) < 3 {
		return "", fmt.Errorf("Attendance week data not in correct form: %v", week)
	}

	return fmt.Sprintf("%v/%v/%v", digits[2], digits[1], digits[0]), nil
}

// LookupDate lookups the id number of the date, and returns its name
func (db Database) LookupDate(id string) (string, error) {

	d, exists := db.dateMap[id]
	if !exists {
		return "", fmt.Errorf("Date not found with id: %v", id)
	}
	return d, nil
}

// LookupResultset looks up the id number of the resultset and returns
// its name
func (db Database) LookupResultset(id string) (string, error) {

	rs, exists := db.resultsetMap[id]
	if !exists {
		return "", fmt.Errorf("Resultset not found with id: %v", id)
	}
	return rs, nil
}

// LookupNatYear looks up the id number of the National Dataset and returns
// its name
func (db Database) LookupNatYear(id string) (string, error) {

	ny, exists := db.natYearMap[id]
	if !exists {
		return "", fmt.Errorf("National data not found with id: %v", id)
	}
	return ny, nil
}
