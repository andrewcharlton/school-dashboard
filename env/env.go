// Package env pulls together the environment variables
// needed for the dashboard application:
//
// * Database
// * Templates
package env

import (
	"html/template"

	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/andrewcharlton/school-dashboard/templates"
)

// Env is a container for environment variables
type Env struct {
	database.Database
	Templates *template.Template
}

// Connect to the database, and load up template files
func Connect(filename string) (Env, error) {

	e := Env{}
	var err error

	e.Database, err = database.Connect(filename)
	if err != nil {
		return Env{}, err
	}

	e.Templates, err = templates.Parse()
	if err != nil {
		return Env{}, err
	}

	return e, nil
}
