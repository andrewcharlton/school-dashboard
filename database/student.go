package database

import (
	"database/sql"
	"errors"
	"strings"
)

// A StudentFilter holds all the fields necessary to uniquely
// identify a student and the grades etc.
type StudentFilter struct {
	UPN       string
	Date      string
	Resultset string
	NatYear   string
}

// Student creates a student object based on the
func (db sqliteDB) Student(f StudentFilter) (analysis.Student, error) {



func (db Database) studentDetails(f StudentFilter) (student.Student, error) {

	s := student.Student{}
	sen := student.SENInfo{}
	ks2 := student.KS2Info{}

	// Load student details
	row := db.stmts["student"].QueryRow(f.UPN, f.Date)
	var aps float64
	var band, en, ma, av, re, wr, gps string
	var male bool
	err := row.Scan(&s.UPN, &s.Surname, &s.Forename, &s.Year, &s.Form,
		&s.PP, &s.EAL, &male, &s.Ethnicity, &sen.Status, &sen.Need,
		&sen.Info, &sen.Strategies, &sen.Access, &sen.IEP, &aps,
		&band, &en, &ma, &av, &re, &wr, &gps)
	if err == sql.ErrNoRows {
		return analysis.Student{}, errors.New("Student not on roll at this date")
	}
	if err != nil {
		return analysis.Student{}, err
	}
	s.SEN = sen
	s.KS2 = analysis.KS2Info{Exists: aps > 0.0, APS: aps, Band: band,
		En: strings.ToUpper(en), Ma: strings.ToUpper(ma),
		Av: strings.ToUpper(av), Re: strings.ToUpper(re),
		Wr: strings.ToUpper(wr), GPS: strings.ToUpper(gps)}
	if male {
		s.Gender = "Male"
	} else {
		s.Gender = "Female"
	}
	return s, nil
}


func (db Database) studentCourses(f StudentFilter) (map[string]subject.Course, error) {

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
	defer rows.Close()

	s.Courses = map[string]analysis.Course{}
	for rows.Next() {
		var subjID, effort int
		var subj, grade string
		err := rows.Scan(&subjID, &subj, &grade, &effort)
		if err != nil {
			return analysis.Student{}, err
		}

		// Don't overwrite existing course - ensures first/best exam is used.
		_, exists := s.Courses[subj]
		if exists {
			continue
		}

		subject := db.subjects[subjID]
		c := subject.Course{subject, subject.Gradeset[grade], effort, "", ""}
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
