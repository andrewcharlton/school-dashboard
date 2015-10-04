package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/database"
)

func (e Env) StudentHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	fmt.Fprintf(w, r.URL.RawQuery)
	upn, exists := query["UPN"]
	if !exists {
		fmt.Fprintf(w, "Try searching for someone!")
		return
	}

	f := database.StudentFilter{upn[0], "1", "", ""}
	s, err := e.DB.LoadStudent(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, s.Surname+", "+s.Forename)
}
