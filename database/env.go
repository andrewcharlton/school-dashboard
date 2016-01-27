package database

import (
	"fmt"
	"html/template"

	"github.com/andrewcharlton/school-dashboard/national"
	"github.com/andrewcharlton/school-dashboard/templates"
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

// LoadTemplates parses the contents of the templates package
// and creates an html.Template from the contents
func (e *Env) LoadTemplates() error {

	var t *template.Template
	for name, contents := range templates.Templates {
		var tmpl *template.Template

		if t == nil {
			t = template.New(name)
		}

		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}

		_, err := tmpl.Parse(contents)
		if err != nil {
			return err
		}
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
	}
	e.Ethnicities = append(e.Ethnicities, "Other")

	return nil
}

// LoadNationals result data
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

// LookupDate lookups the id number of the date, and returns its name
func (e Env) LookupDate(id string) (string, error) {

	for _, d := range e.Dates {
		if d.ID == id {
			return d.Name, nil
		}
	}
	return "", fmt.Errorf("Date not found with id: %v", id)
}

// LookupResultset looks up the id number of the resultset and returns
// its name
func (e Env) LookupResultset(id string) (string, error) {

	for _, r := range e.Resultsets {
		if r.ID == id {
			return r.Name, nil
		}
	}
	return "", fmt.Errorf("Resultset not found with id: %v", id)
}

// LookupNatYear looks up the id number of the National Dataset and returns
// its name
func (e Env) LookupNatYear(id string) (string, error) {

	for _, n := range e.NatYears {
		if n.ID == id {
			return n.Name, nil
		}
	}
	return "", fmt.Errorf("National data not found with id: %v", id)
}
