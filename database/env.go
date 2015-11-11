package database

import (
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/andrewcharlton/school-dashboard/national"
)

// An Env contains all of the environment variables
type Env struct {
	// Database connection
	DB Database

	// Misc config variables
	Config Config

	// HTML templates
	Templates *template.Template

	// Cached lists for menu/filter items.
	Dates       []Lookup
	Resultsets  []Lookup
	Ethnicities []string
	// OtherEths tells us whether an ethnicity should
	// be collapsed into the 'Other' category.
	OtherEths map[string]bool

	// National Data for each year
	NatYears  []Lookup
	Nationals map[string]national.National
}

// Connect to the database and initialise all
// environment variables.
func Connect(filename string) (Env, error) {

	db, err := ConnectSQLite(filename)
	if err != nil {
		return Env{}, err
	}

	e := Env{DB: db}
	if err = e.LoadConfig(); err != nil {
		return Env{}, err
	}

	if err = e.LoadTemplates(); err != nil {
		return Env{}, err
	}

	if err = e.LoadFilterItems(); err != nil {
		return Env{}, err
	}

	if err = e.LoadNationals(); err != nil {
		return Env{}, err
	}

	return e, nil
}

// LoadConfig loads up all the config variables into
// memory.
func (e *Env) LoadConfig() error {

	cfg, err := e.DB.Config()
	if err != nil {
		return err
	}

	e.Config = cfg
	return nil
}

// LoadTemplates searches for all template files in
// ./templates and parses them.
func (e *Env) LoadTemplates() error {

	filenames := []string{}
	files, err := ioutil.ReadDir("templates/")
	if err != nil {
		return err
	}

	// Check they are .html and give full path back
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

// LoadFilterItems populates the lists for the
func (e *Env) LoadFilterItems() error {

	dates, err := e.DB.Dates()
	if err != nil {
		return err
	}
	e.Dates = dates

	rs, err := e.DB.Resultsets()
	if err != nil {
		return err
	}
	e.Resultsets = rs

	ethnicities, err := e.DB.Ethnicities()
	if err != nil {
		return err
	}

	// Populate the first 8 groups
	e.Ethnicities = []string{}
	e.OtherEths = map[string]bool{}
	for n, eth := range ethnicities {
		if n < 8 {
			e.Ethnicities = append(e.Ethnicities, eth.Name)
			e.OtherEths[eth.Name] = false
		} else {
			e.OtherEths[eth.Name] = true
		}
		e.Ethnicities = append(e.Ethnicities, "Other")
	}

	return nil
}

// LoadNational result data
func (e *Env) LoadNationals() error {

	years, err := e.DB.NationalYears()
	if err != nil {
		return err
	}
	e.NatYears = years

	e.Nationals = map[string]national.National{}
	for _, y := range years {
		nat, err := e.DB.National(y.ID)
		if err != nil {
			return err
		}
		e.Nationals[y.ID] = nat
	}

	return nil
}
