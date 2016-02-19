package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/analysis/group"
	"github.com/andrewcharlton/school-dashboard/env"
)

// Search returns a student search page
func Search(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println(r.URL.RawPath, r.URL.RawQuery)

		header(e, w, r, 0)
		defer footer(e, w, r)

		name := r.URL.Query().Get("name")

		f := getFilter(e, r)
		g, err := e.Search(name, f)

		data := struct {
			Query template.URL
			Name  string
			Group group.Group
		}{
			template.URL(r.URL.RawQuery),
			name,
			g,
		}

		err = e.Templates.ExecuteTemplate(w, "studentsearch.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}
