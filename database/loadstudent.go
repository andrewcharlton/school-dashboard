package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" //database driver

	"github.com/andrewcharlton/school-dashboard/analysis/student"
	"github.com/andrewcharlton/school-dashboard/analysis/subject"
)

// Student loads the details of an individual student
func (db Database) Student(upn string, f Filter) (student.Student, error) {

	s, err := db.loadStudent(upn, f)
	if err != nil {
		return student.Student{}, err
	}

	// Set Attainment 8 data - based on KS2.APS
	att8, exists := db.Attainment8[f.NatYear][selectProgress8(s.KS2.APS)]
	s.SetNationals(f.NatYear, att8, exists)

	err = db.loadStudentResults(&s, f)
	if err != nil {
		return student.Student{}, err
	}

	err = db.loadStudentClasses(&s, f)
	if err != nil {
		return student.Student{}, err
	}

	err = db.loadStudentAttendance(&s, f)
	if err != nil {
		return student.Student{}, err
	}

	return s, nil
}

// loadStudent loads up the basic student details and creates a Student object from them
func (db Database) loadStudent(upn string, f Filter) (student.Student, error) {

	row := db.stmts["student"].QueryRow(upn, f.Date)

	s := student.Student{}
	sen := student.SENInfo{}
	ks2 := student.KS2Info{}
	err := row.Scan(&s.UPN, &s.Surname, &s.Forename, &s.Year, &s.Form,
		&s.PP, &s.EAL, &s.Gender, &s.Ethnicity, &sen.Status, &sen.Need,
		&sen.Info, &sen.Strategies, &sen.Access, &sen.IEP, &ks2.APS,
		&ks2.Band, &ks2.En, &ks2.Ma, &ks2.Av, &ks2.Re, &ks2.Wr, &ks2.GPS)
	switch {
	case err == sql.ErrNoRows:
		return student.Student{}, fmt.Errorf("Student not on roll at this time.")
	case err != nil:
		return student.Student{}, err
	}
	s.SEN = sen
	s.KS2 = ks2

	return s, nil
}

// loadStudentResults populates the lists of subjects where the student
// has a grade and/or effort from the resultset specified in the filter.
func (db Database) loadStudentResults(s *student.Student, f Filter) error {

	var rows *sql.Rows
	var err error
	switch f.Resultset {
	case "1":
		rows, err = db.stmts["firstExams"].Query(s.UPN)
	case "2":
		rows, err = db.stmts["bestExams"].Query(s.UPN)
	default:
		rows, err = db.stmts["results"].Query(s.UPN, f.Resultset)
	}

	if err != nil {
		return err
	}
	defer rows.Close()

	s.Results = map[string]subject.Result{}
	for rows.Next() {
		var subjID, effort int
		var subjName, grade string
		err := rows.Scan(&subjID, &subjName, &grade, &effort)
		if err != nil {
			return err
		}

		// Don't overwrite existing courses - this, plus the ordering
		// for exams means only first/best will be included
		_, exists := s.Results[subjName]
		if exists {
			continue
		}

		subj := db.Subjects[subjID]
		r := subject.Result{subj, subj.Gradeset[grade], effort, "", ""}
		s.Results[subjName] = r
	}

	return nil
}

// loadStudentClasses populates the classnames and teachers for any subjects
// which are present in a student's list of results.
func (db Database) loadStudentClasses(s *student.Student, f Filter) error {

	rows, err := db.stmts["classes"].Query(s.UPN, f.Date)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var subject, class, teacher string
		err := rows.Scan(&subject, &class, &teacher)
		if err != nil {
			return err
		}

		r, exists := s.Results[subject]
		if exists {
			r.Class = class
			r.Teacher = teacher
			s.Results[subject] = r
		}
	}

	return nil
}

// loadStudentAttendance populates the latest attendance figures for a
// student
func (db Database) loadStudentAttendance(s *student.Student, f Filter) error {

	row := db.stmts["attendance"].QueryRow(s.UPN, f.Date)
	att := student.AttendanceInfo{}
	err := row.Scan(&att.Week, &att.Possible, &att.Absences,
		&att.Unauthorised, &att.Sessions[0], &att.Sessions[1],
		&att.Sessions[2], &att.Sessions[3], &att.Sessions[4],
		&att.Sessions[5], &att.Sessions[6], &att.Sessions[7],
		&att.Sessions[8], &att.Sessions[9])

	switch {
	case err == sql.ErrNoRows || err == nil:
		s.Attendance = att
		return nil
	case err != nil:
		return err
	}

	return nil
}

// Select the correct progress 8 set of national data to use for a student.
func selectProgress8(aps float64) string {

	switch {
	case aps < 0.5:
		return ""
	case aps < 1.55*6:
		return "1.5"
	case aps < 2.05*6:
		return "2.0"
	case aps < 2.55*6:
		return "2.5"
	case aps < 2.85*6:
		return "2.8"
	case aps >= 5.75*6:
		return "5.8"
	default:
		return fmt.Sprintf("%1.1f", aps/6)
	}
}
