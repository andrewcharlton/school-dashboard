package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/andrewcharlton/school-dashboard/analysis/group"
)

// A NewsItem contains the details of a news item
type NewsItem struct {
	Date    string
	Comment string
}

// News returns a list of news items
func (db Database) News() ([]NewsItem, error) {

	news := []NewsItem{}
	rows, err := db.stmts["news"].Query()
	if err != nil {
		return []NewsItem{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var date, comm string
		err := rows.Scan(&date, &comm)
		if err != nil {
			return []NewsItem{}, err
		}

		split := strings.Split(date, "-")
		if len(split) == 3 {
			formattedDate := fmt.Sprintf("%v-%v-%v", split[2], split[1], split[0])
			news = append(news, NewsItem{formattedDate, comm})
		}
	}

	return news, nil
}

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

// Classes returns a list of classes that exists for a subject,
// at a particular date.
func (db Database) Classes(subjID, dateID string) ([]string, error) {

	rows, err := db.stmts["classlist"].Query(dateID, subjID)
	if err == sql.ErrNoRows {
		return []string{}, fmt.Errorf("No classes attached to this subject.")
	}
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	classes := []string{}
	for rows.Next() {
		var class string
		err := rows.Scan(&class)
		if err != nil {
			return []string{}, err
		}
		classes = append(classes, class)
	}

	return classes, nil
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

// Search returns a list of students who's names match the search string
func (db Database) Search(name string, f Filter) (group.Group, error) {

	name = "%" + strings.Replace(name, "*", "%", -1) + "%"
	rows, err := db.stmts["search"].Query(f.Date, name, name)
	if err == sql.ErrNoRows {
		return group.Group{}, fmt.Errorf("No students called %v found", name)
	}
	if err != nil {
		return group.Group{}, err
	}
	defer rows.Close()

	upns := []string{}
	for rows.Next() {
		var upn string
		err := rows.Scan(&upn)
		if err != nil {
			return group.Group{}, err
		}
		upns = append(upns, upn)
	}

	return db.basicGroup(upns, f)
}
