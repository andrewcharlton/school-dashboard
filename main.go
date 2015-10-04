package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/database"
)

// Env contains all of the
type Env struct {
	DB database.SchoolDB
}

func LoadEnvironment(filename string) (Env, error) {

	db, err := database.Connect(filename)
	if err != nil {
		return Env{}, err
	}

	e := Env{db}
	return e, nil
}

func main() {

	// Get filename
	args := flag.Args()
	filename := "school.db"
	if len(args) >= 1 {
		filename = args[1]
	}

	// Load environment
	e, err := LoadEnvironment(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer e.DB.Close()

	// Handlers
	http.HandleFunc("/", e.Index)
	http.HandleFunc("/students/", e.StudentHandler)

	// Start serer
	http.ListenAndServe(":8080", nil)
}

func (e Env) Index(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello")

}
