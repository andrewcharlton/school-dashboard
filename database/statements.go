package database

import "database/sql"

var query_sql = map[string]string{

	"search_student": `SELECT upn
						FROM students
						WHERE ((forename || " " || surname) LIKE "%?%")
						OR ((surname || ", " || forename) LIKE "%?%")`,

	"load_student": `SELECT upn, surname, forename, year, form,
						pp, eal, gender, ethnicity, sen_status,
						sen_info, sen_strat, ks2_aps, ks2_band,
						ks2_en, ks2_ma, ks2_av
						FROM students
						WHERE upn=? AND date=?`,
}

// closeQueries closes all prepared statements
func (db *SchoolDB) closeQueries() error {

	for _, s := range db.queries {
		err := s.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// prepQueries prepares statements from all query strings
// in query_sql
func (db *SchoolDB) prepQueries() error {

	db.queries = map[string]*sql.Stmt{}
	for key, query := range query_sql {
		s, err := db.DB.Prepare(query)
		if err != nil {
			return err
		}
		db.queries[key] = s
	}

	return nil
}
