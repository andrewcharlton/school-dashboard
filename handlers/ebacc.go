package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/analysis/group"
	"github.com/andrewcharlton/school-dashboard/env"
)

// EBacc produces a page with summary figures of how the group
// has achieved in the EBacc.
func EBacc(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if redir := checkRedirect(e, w, r, 2); redir {
			return
		}
		header(e, w, r, 2)
		defer footer(e, w, r)

		f := getFilter(e, r)
		g, err := e.GroupByFilter(f)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		groups := [3]group.Group{}
		for _, s := range g.Students {
			switch {
			case s.EBacc().Achieved:
				groups[0].Students = append(groups[0].Students, s)
			case s.EBacc().Entered:
				groups[1].Students = append(groups[1].Students, s)
			default:
				groups[2].Students = append(groups[2].Students, s)
			}
		}

		data := struct {
			Query        template.URL
			Areas        []string
			Group        group.Group
			SubGroups    [3]group.Group
			GroupHeaders []string
		}{
			template.URL(r.URL.RawQuery),
			[]string{"E", "M", "S", "H", "L"},
			g,
			groups,
			[]string{"Achieving EBacc", "Eligible for EBacc", "Not Eligible for EBacc"},
		}

		err = e.Templates.ExecuteTemplate(w, "ebacc.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}

}
