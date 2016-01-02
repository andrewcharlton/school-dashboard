package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/database"
)

type effort struct {
	UPN     string
	Name    string
	Scores  map[int]int
	Average float64
	Prog8   float64
}

func Effort(e database.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if redir := checkRedirect(e, queryOpts{true, true}, w, r); redir {
			return
		}

		Header(e, w, r)
		FilterPage(e, w, r, false)
		defer Footer(e, w, r)

		f := GetFilter(e, r)
		g, err := e.DB.GroupByFilter(f)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}

		nat := e.Nationals[f.NatYear]

		efforts := []effort{}
		for _, s := range g.Students {
			eff := effort{UPN: s.UPN, Name: s.Name(), Scores: map[int]int{}}
			total, num := 0, 0
			for _, c := range s.Courses {
				eff.Scores[c.Effort] += 1
				total += c.Effort
				num += 1
			}
			if num == 0 {
				eff.Average = float64(0)
			} else {
				eff.Average = float64(total) / float64(num)
			}

			b := s.Basket()
			natP8, err := nat.Progress8(s.KS2.APS)
			if err == nil {
				eff.Prog8 = b.Progress8(natP8).Pts
			}

			efforts = append(efforts, eff)
		}

		data := struct {
			Efforts []effort
			Query   template.URL
		}{
			efforts,
			template.URL(r.URL.RawQuery),
		}

		e.Templates.ExecuteTemplate(w, "effort.tmpl", data)
	}
}
