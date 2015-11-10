package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/env"
)

func EnglishAndMaths(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		Header(e, w, r)
		FilterPage(e, w, r, false)
		defer Footer(e, w, r)

		f := GetFilter(e, r)
		g, err := e.DB.GroupByFilter(f)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}
		nat := e.Nationals[f.NatYear]

		type student struct {
			UPN    string
			Name   string
			EnGrd  string
			EnEff  int
			MaGrd  string
			MaEff  int
			Basics bool
			AvEff  float64
			P8     float64
		}

		data := struct {
			Students    []student
			Query       template.URL
			Cohort      int
			EnPass      int
			MaPass      int
			BothPass    int
			EnPassPct   float64
			MaPassPct   float64
			BothPassPct float64
		}{
			Students: []student{},
			Query:    template.URL(ShortenQuery(e, r.URL.Query())),
		}

		for _, s := range g.Students {
			data.Cohort += 1
			if s.Basics().AchB {
				data.BothPass += 1
			}

			en, exists := s.Courses["English"]
			var enGrd string
			var enEff int
			if exists {
				enGrd = en.Grd
				enEff = en.Effort
				if en.L2Pass {
					data.EnPass += 1
				}
			}

			ma, exists := s.Courses["Mathematics"]
			var maGrd string
			var maEff int
			if exists {
				maGrd = ma.Grd
				maEff = ma.Effort
				if ma.L2Pass {
					data.MaPass += 1
				}
			}

			data.Students = append(data.Students, student{s.UPN,
				s.Name(),
				enGrd,
				enEff,
				maGrd,
				maEff,
				s.Basics().AchB,
				s.Effort().Pts,
				s.Basket().Progress8(nat).Pts,
			})
		}

		if data.Cohort > 0 {
			data.EnPassPct = float64(100) * float64(data.EnPass) / float64(data.Cohort)
			data.MaPassPct = float64(100) * float64(data.MaPass) / float64(data.Cohort)
			data.BothPassPct = float64(100) * float64(data.BothPass) / float64(data.Cohort)
		}

		err = e.Templates.ExecuteTemplate(w, "em.tmpl", data)
		if err != nil {
			fmt.Println(err)
		}
	}
}
