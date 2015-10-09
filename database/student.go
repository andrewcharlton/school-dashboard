package database

import (
	"github.com/andrewcharlton/school-dashboard/analysis"
)

// A StudentFilter provides all of the fields necessary to
// uniquely identify a student, and the classes they are
// in/grades they have achieved.
type StudentFilter struct {
	UPN       string
	Date      string
	Resultset string
	Exams     string
}

// LoadStudent retrieves a student's information from the database
// and returns a Student object holding that info.
func (db SchoolDB) LoadStudent(f StudentFilter) (analysis.Student, error) {

	query := `SELECT upn, surname, forename, year, form,
			  pp, eal, gender, ethnicity, sen_status, sen_info,
			  sen_strat, ks2_aps, ks2_band, ks2_en, ks2_ma, ks2_av
			  FROM students
			  WHERE upn=? AND date=?`

	row := db.DB.QueryRow(query, f.UPN, f.Date)

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
	return s, nil
}
