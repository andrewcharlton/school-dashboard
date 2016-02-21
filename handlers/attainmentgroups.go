package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/env"
)

func AttainmentGroups(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if redir := checkRedirect(e, w, r, 1); redir {
			return
		}
		header(e, w, r, 1)
		defer footer(e, w, r)

		f := getFilter(e, r)
		g, err := e.GroupByFilter(f)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		data := struct {
			Query  template.URL
			Groups []subGroup
		}{
			template.URL(r.URL.RawQuery),
			subGroups(g),
		}

		err = e.Templates.ExecuteTemplate(w, "attainmentgroups.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}
