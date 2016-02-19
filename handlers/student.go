package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/andrewcharlton/school-dashboard/analysis/student"
	"github.com/andrewcharlton/school-dashboard/env"
)

func Student(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if redir := checkRedirect(e, w, r, 0); redir {
			return
		}
		header(e, w, r, 0)
		defer footer(e, w, r)

		path := strings.Split(r.URL.Path, "/")
		if len(path) < 3 {
			fmt.Fprintf(w, "Error: Invalid path")
			return
		}
		upn := path[2]

		f := getFilter(e, r)
		s, err := e.Student(upn, f)
		if err != nil {
			fmt.Fprintf(w, "Error: %", err)
		}

		data := struct {
			Student student.Student
		}{
			s,
		}

		err = e.Templates.ExecuteTemplate(w, "student.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}
