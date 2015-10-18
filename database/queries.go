package database

import "database/sql"

var query_sql = map[string]string{

	"search_student": `SELECT surname, forename
						FROM students
						WHERE ((forename || " " || surname) LIKE "%j%")
						OR ((surname || ", " || forename) LIKE "%j%")`,
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
