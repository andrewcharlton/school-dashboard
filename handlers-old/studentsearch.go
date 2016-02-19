package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/andrewcharlton/school-dashboard/env"
)

// Search returns a student search page
func Search(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if redir := checkRedirect(e, w, r, 0); redir {
			return
		}
		header(e, w, r, 0)
		defer footer(e, w, r)

		name := ""
		path := strings.Split(r.URL.Path, "/")
		for i := len(path) - 1; i > 1; i-- {
			if path[i] != "" {
				name = path[i]
				break
			}
		}

		f := GetFilter(e, r)
		list, err := e.Search(name, f.Date)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		data := struct {
			Name     string
			Query    template.URL
			Results  bool
			Students []database.StudentLookup
		}{
			name,
			template.URL(r.URL.RawQuery),
			(len(list) > 0),
			list,
		}

		err = e.Templates.ExecuteTemplate(w, "studentsearch.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}
