package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/andrewcharlton/school-dashboard/database"
)

// Env contains all of the
type Env struct {
	DB        database.Database
	Templates *template.Template

	// Data required for filter drop-down
	Dates       []Option
	Resultsets  []Option
	Ethnicities []string
	OtherEths   map[string]bool
}

// Connect to the database and return an Environment variable
// containing the connection.
func Connect(filename string) (Env, error) {

	db, err := database.ConnectSQLite("school.db")
	if err != nil {
		return Env{}, err
	}

	e := Env{DB: db, Templates: nil, Dates: []Option{},
		Resultsets: []Option{}, Ethnicities: []string{},
		OtherEths: map[string]bool{}}
	return e, nil
}

// LoadTemplates searches for all template files in
// ./templates and parses them.
func (e *Env) LoadTemplates() error {

	filenames := []string{}
	files, err := ioutil.ReadDir("templates/")
	if err != nil {
		return err
	}

	// Check they are .tmpl and give full path back
	for _, f := range files {
		if strings.Contains(f.Name(), ".tmpl") {
			filenames = append(filenames, "templates/"+f.Name())
		}
	}

	t, err := template.ParseFiles(filenames...)
	if err != nil {
		return err
	}

	e.Templates = t
	return nil
}

// LoadDates loads all effective dates that are marked to
// be listed.
func (e *Env) LoadDates() error {
	rows, err := e.DB.Query(`SELECT id, date FROM dates
							WHERE list=1
							ORDER BY date DESC`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, date string
		err := rows.Scan(&id, &date)
		if err != nil {
			return err
		}
		e.Dates = append(e.Dates, Option{id, date})
	}

	return nil
}

// LoadResultsets Loads
func (e *Env) LoadResultsets() error {

	e.Resultsets = []Option{{"0", "Exams Only"}}
	rows, err := e.DB.Query(`SELECT id, resultset FROM resultsets
							WHERE is_exam=0 AND list=1
							ORDER BY date DESC`)
	if err != nil {
		return err
	}

	for rows.Next() {
		var id, rs string
		err := rows.Scan(&id, &rs)
		if err != nil {
			return err
		}
		e.Resultsets = append(e.Resultsets, Option{id, rs})
	}

	return nil
}

// LoadEthnicities loads up the most common ethnicities
// for the filter drop down.  For ethnicities that appear
// rarely, these are folded into 'Other' and tagged in
// the OtherEths folder.
func (e *Env) LoadEthnicities() error {

	e.OtherEths = map[string]bool{}
	rows, err := e.DB.Query(`SELECT ethnicity, COUNT(1) as n
							FROM students
							GROUP BY ethnicity
							ORDER BY n DESC`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var eth string
		var count int
		err := rows.Scan(&eth, &count)
		if err != nil {
			return err
		}
		e.Ethnicities = append(e.Ethnicities, eth)
		if count < 10 {
			e.OtherEths[eth] = true
		}
	}

	return nil

}

// LoadEnvironment wraps the loading of templates, dates
// resultsets and ethnicities.
func (e *Env) LoadEnvironment() error {

	err := e.LoadTemplates()
	if err != nil {
		return err
	}

	err = e.LoadDates()
	if err != nil {
		return err
	}

	err = e.LoadResultsets()
	if err != nil {
		return err
	}

	err = e.LoadEthnicities()
	if err != nil {
		return err
	}

	return nil
}

func main() {

	// Get filename
	args := flag.Args()
	filename := "school.db"
	if len(args) >= 1 {
		filename = args[1]
	}

	// Connect to database
	e, err := Connect(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer e.DB.Close()

	// Load Environment Variables
	err = e.LoadEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// Static fileserver
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handlers
	http.HandleFunc("/", e.Index)
	http.HandleFunc("/students/", e.StudentHandler)

	// Start serer
	http.ListenAndServe(":8080", nil)
}

func (e Env) Index(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	f, err := e.GetFilter(query)

	if err != nil {
		log.Fatal(err)
	}

	err = e.Header(w, r)
	if err != nil {
		log.Fatal(err)
	}

	err = e.FilterPage(w, f, false)
	if err != nil {
		log.Fatal(err)
	}

	err = e.Footer(w)
	if err != nil {
		log.Fatal(err)
	}

}
