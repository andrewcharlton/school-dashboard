package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/analysis/group"
	"github.com/andrewcharlton/school-dashboard/env"
)

type emList struct {
	Name        string
	Percentages []float64
}

func EnglishAndMaths(e env.Env) http.HandlerFunc {
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

		groups := []group.Group{{}, {}, {}, {}}
		for _, s := range g.Students {
			eng := s.EBaccArea("E").Achieved
			maths := s.EBaccArea("M").Achieved
			switch {
			case eng && maths:
				groups[2].Students = append(groups[2].Students, s)
			case eng:
				groups[0].Students = append(groups[0].Students, s)
			case maths:
				groups[1].Students = append(groups[1].Students, s)
			default:
				groups[3].Students = append(groups[3].Students, s)
			}
		}

		pcts := []float64{}
		for _, grp := range groups {
			pcts = append(pcts, float64(len(grp.Students))/float64(len(g.Students)))
		}

		data := struct {
			Query  template.URL
			Names  []string
			Groups []group.Group
			Pcts   []float64
		}{
			template.URL(r.URL.RawQuery),
			[]string{"English Only", "Mathematics Only", "English & Maths", "Neither"},
			groups,
			pcts,
		}

		err = e.Templates.ExecuteTemplate(w, "em.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}

}
