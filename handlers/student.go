package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/andrewcharlton/school-dashboard/env"
)

func Student(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		Header(e, w, r)
		FilterPage(e, w, r, true)

		upn := ""
		path := strings.Split(r.URL.Path, "/")
		for i := len(path) - 1; i > 1; i-- {
			if path[i] != "" {
				upn = path[i]
				break
			}
		}

		query := r.URL.Query()
		f := database.StudentFilter{
			upn,
			query.Get("date"),
			query.Get("resultset"),
		}

		s, err := e.DB.Student(f)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		fmt.Fprintf(w, "Student: %v, %v", s.Surname, s.Forename)
		for _, c := range s.Courses {
			fmt.Fprintf(w, "<br>Results: %v, %v", c.Subj, c.Grd)
		}

		Footer(e, w, r)
	}
}