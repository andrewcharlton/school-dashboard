// Package env contains functions to
// construct and initialise the environment
// variable for the dashboard application
package env

import (
	"github.com/andrewcharlton/school-dashboard/database"
)

// An Env contains all of the environment variables
// such as the database connection.
type Env struct {
	DB database.Database
}

// Connect to the database and initialise all
// environment variables.
func Connect(filename string) (Env, error) {

	db, err := database.ConnectSQLite(filename)
	if err != nil {
		return Env{}, err
	}

	return Env{DB: db}, nil
}
