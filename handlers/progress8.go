package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
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
		g, err := e.GroupByFilter(f)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		nat, exists := e.Attainment8[f.NatYear]
		if !exists {
			fmt.Fprintf(w, "Error: %v", err)
		}

		natLine := points{}
		for ks2, att8 := range nat {
			n, err := strconv.ParseFloat(ks2, 64)
			if err != nil {
				fmt.Fprintf(w, "Error: %v", err)
			}
			natLine = append(natLine, point{X: n, Y: att8.Overall})
		}
		sort.Sort(natLine)

		pupilData := [4]points{}
		for _, s := range g.Students {

			p8 := s.Basket().Overall()
			var key int
			switch {
			case s.PP && s.Gender == 1:
				key = 0
			case s.Gender == 1:
				key = 1
			case s.PP && s.Gender == 0:
				key = 2
			case s.Gender == 0:
				key = 3
			}

			pupilData[key] = append(pupilData[key], point{X: s.KS2.APS / 6,
				Y:    p8.Attainment,
				Name: s.Name(),
				P8:   p8})
		}

		data := struct {
			Query     template.URL
			Group     group.Group
			NatLine   points
			PupilData [4]points
		}{
			template.URL(r.URL.RawQuery),
			g,
			natLine,
			pupilData,
		}

		err = e.Templates.ExecuteTemplate(w, "progress8.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}

// Progress8Groups calculates the progress 8 scores for each group of students
func Progress8Groups(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if redir := checkRedirect(e, queryOpts{true, false}, w, r); redir {
			return
		}

		Header(e, w, r)
		FilterPage(e, w, r, false)
		defer Footer(e, w, r)

		f := GetFilter(e, r)
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

		err = e.Templates.ExecuteTemplate(w, "progress8groups.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}
