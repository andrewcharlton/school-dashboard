package db

import (
	"sort"
)

// A lookupRow holds the id and name of
// a field lookup
type lookupRow struct {
	ID
	Name
}

// A Lookup holds a set of ID/Name pairs
type Lookup []lookupRow

func (l Lookup) Len() int      { return Len(l) }
func (l Lookup) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l Lookup) Less(i, j int) { return l[i].ID < l[j].ID }

// Dates returns a sorted list of all Dates in the database that
// are marked to be listed.
func (db Database) Dates() (Lookup, error) {

	rows, err := db.db.Query("SELECT id, date FROM dates WHERE list=1")
	if err != nil {
		return Lookup{}, err
	}
	defer rows.Close()

	dates := Lookup{}
	for rows.Next() {
		var r lookupRow
		err := rows.Scan(&r.ID, &r.Name)
		if err != nil {
			return Lookup{}, err
		}
		dates = append(dates, r)
	}

	sort.Sort(dates)
	return dates, nil
}

// Resultsets returns a sorted list of all Resultsets in the database.
// An 'Exams Only' option encapsulates all individual exam resultsets,
// all other resultsets marked to be listed are included.
func (db Database) Resultsets() (Lookup, error) {

	rows, err := db.db.Query(`SELECT id, resultset FROM resultsets
								WHERE is_exam=0 AND list=1`)
	if err != nil {
		return Lookup{}, err
	}
	defer rows.Close()

	rs := Lookup{{0, "Exams Only"}}
	for rows.Next() {
		var r lookupRow
		err := rows.Scan(&r.ID, &r.Name)
		if err != nil {
			return Lookup{}, err
		}
		rs = append(rs, r)
	}

	sort.Sort(rs)
	return rs, nil
}

type ethnicity struct {
	Ethnicity string
	Count     int
}

// Ethnicities contains a list of all ethnicties and the number
// of students who have that ethnicity.
type Ethnicities []ethnicity

func (e Ethnicities) Len() int      { return Len(e) }
func (e Ethnicities) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e Ethnicities) Less(i, j int) { return e[i].Count < e[j].Count }

// Ethnicities returns all the distinct ethnicities present in
// the database, and the frequency that each appears with.
func (db Database) Ethnicities() (Ethnicities, error) {

	rows, err := db.db.Query(`SELECT ethnicity, COUNT(1) as n
								FROM students
								GROUP BY ethnicity`)
	if err != nil {
		return Ethnicities{}, err
	}
	defer rows.Close()

	eth := Ethnicities{}
	for rows.Next() {
		var e ethnicity
		err := rows.Scan(&e.Ethnicity, &e.Count)
		if err != nil {
			return Ethnicities{}, err
		}
		eth = append(eth, e)
	}

	sort.Sort(sort.Reverse(eth))
	return eth, nil
}
