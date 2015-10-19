package db

// A Lookup holds an ID/Name pair for an item in the database.
type Lookup struct {
	ID   string
	Name string
}

// Dates returns a sorted list of all Dates in the database that
// are marked to be listed.
func (db Database) Dates() ([]Lookup, error) {

	rows, err := db.db.Query(`SELECT id, date FROM dates WHERE list=1
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

// Resultsets returns a sorted list of all Resultsets in the database.
// An 'Exams Only' option encapsulates all individual exam resultsets,
// all other resultsets marked to be listed are included.
func (db Database) Resultsets() ([]Lookup, error) {

	rows, err := db.db.Query(`SELECT id, resultset FROM resultsets
								WHERE is_exam=0 AND list=1
								ORDER BY id`)
	if err != nil {
		return []Lookup{}, err
	}
	defer rows.Close()

	rs := []Lookup{{0, "Exams Only"}}
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

// An Ethnicity holds the name of an Ethnic group, as well as the
// number of students in that group.
type Ethnicity struct {
	Name  string
	Count int
}

// Ethnicities returns all the distinct ethnicities present in
// the database, and the frequency that each appears with.
// Only students who are present in listed dates are counted.
func (db Database) Ethnicities() ([]Ethnicity, error) {

	rows, err := db.db.Query(`SELECT ethnicity, COUNT(1) as n
								FROM students
								GROUP BY ethnicity`)
	if err != nil {
		return []Ethnicity{}, err
	}
	defer rows.Close()

	eth := []Ethnicity{}
	for rows.Next() {
		var e ethnicity
		err := rows.Scan(&e.Name, &e.Count)
		if err != nil {
			return []Ethncity{}, err
		}
		eth = append(eth, e)
	}

	return eth, nil
}
