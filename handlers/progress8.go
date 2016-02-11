package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/andrewcharlton/school-dashboard/analysis/group"
	"github.com/andrewcharlton/school-dashboard/analysis/student"
	"github.com/andrewcharlton/school-dashboard/env"
)

type point struct {
	X    float64
	Y    float64
	Name string
	P8   student.Progress8Score
}

type points []point

func (p points) Len() int           { return len(p) }
func (p points) Less(i, j int) bool { return p[i].X < p[j].X }
func (p points) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Progress8 returns a handler to produce the Progress 8 page.
func Progress8(e env.Env) http.HandlerFunc {
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
			return
		}

		nat, exists := e.DB.Attainment8[f.NatYear]
		if !exists {
			fmt.Fprintf(w, "Error: %v", err)
		}

		natLine := points{}
		for ks2, pts := range nat {
			n, err := strconv.ParseFloat(ks2, 64)
			if err != nil {
				fmt.Fprintf(w, "Error: %v", err)
			}
			natLine = append(natLine, point{X: n, Y: pts})
		}

		pupilData := map[string]points{}
		for _, s := range g.Students {

			p8 := s.Basket().Overall()
			switch {
			case s.PP && s.Gender == "Male":
				key = "Male - Disadvantaged"
			case s.Gender == "Male":
				key = "Male - Non-Disadvantaged"
			case s.PP && s.Gender == "Female":
				key = "Female - Disadvantaged"
			case s.Gender == "Female":
				key = "Female - Non-Disadvantaged"
			}
			pts := pupilData[int]
			pts = append(pts, point{X: s.KS2.APS / 6,
				Y:    p8.Attainment,
				Name: s.Name(),
				P8:   p8})
			pupilData[int] = pts
		}

		data := struct {
			Filter    template.URL
			Group     group.Group
			NatLine   points
			PupilData map[string]points
			Summary   group.Progress8Summary
		}{
			template.URL(r.URL.RawQuery),
			g,
			natLine,
			pupilData,
			g.Progress8(),
		}

		err = e.Templates.ExecuteTemplate(w, "progress8.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}
