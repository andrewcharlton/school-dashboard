package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/env"
)

type effort struct {
	UPN     string
	Name    string
	Scores  map[int]int
	Average float64
	Prog8   float64
}

func Effort(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		Header(e, w, r)
		FilterPage(e, w, r, false)

		f := GetFilter(e, r)
		g, err := e.DB.GroupByFilter(f)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}

		nat := e.Nationals[f.NatYear]

		efforts := []effort{}
		prog8 := float64(0)
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
			eff.Prog8 = b.Progress8(nat).Pts
			prog8 += eff.Prog8

			efforts = append(efforts, eff)
		}

		data := struct {
			Efforts []effort
			Query   template.URL
			Prog8   float64
		}{
			efforts,
			template.URL(ShortenQuery(e, r.URL.Query())),
			prog8 / float64(len(efforts)),
		}

		e.Templates.ExecuteTemplate(w, "effort.tmpl", data)

		Footer(e, w, r)

	}
}
