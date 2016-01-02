package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"sort"

	"github.com/andrewcharlton/school-dashboard/analysis"
	"github.com/andrewcharlton/school-dashboard/database"
)

func SubjectOverview(e database.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if redir := checkRedirect(e, queryOpts{true, true}, w, r); redir {
			return
		}

		Header(e, w, r)
		FilterPage(e, w, r, false)
		defer Footer(e, w, r)

		subjSummaries, err := subjOverviewData(e, r)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		f := GetFilter(e, r)

		data := struct {
			Subjects []subjSummary
			Year     string
			Query    template.URL
		}{
			subjSummaries,
			f.Year,
			template.URL(r.URL.RawQuery),
		}

		err = e.Templates.ExecuteTemplate(w, "subject-overview.tmpl", data)
		if err != nil {
			fmt.Println(err)
		}
	}
}

type subjData struct {
	Subject *analysis.Subject
	Cohort  int
	KS2     []float64
	PP      []float64
	Points  []float64
	VA      []float64
}

type subjSummary struct {
	Subject *analysis.Subject
	Cohort  int
	HasKS2  bool
	KS2     float64
	PP      float64
	Points  float64
	HasVA   bool
	VA      float64
	AvGrade string
}

func subjOverviewData(e database.Env, r *http.Request) ([]subjSummary, error) {

	f := GetFilter(e, r)
	g, err := e.DB.GroupByFilter(f)
	if err != nil {
		return []subjSummary{}, err
	}

	nat := e.Nationals[f.NatYear]

	data := map[string]subjData{}

	for _, s := range g.Students {
		for subj, c := range s.Courses {
			sd, exists := data[subj]
			if !exists {
				sd = subjData{}
			}

			if sd.Subject == nil {
				sd.Subject = c.Subject
			}
			sd.Cohort++
			if s.KS2.APS != 0 {
				sd.KS2 = append(sd.KS2, s.KS2.APS)
			}
			if s.PP {
				sd.PP = append(sd.PP, 1.0)
			} else {
				sd.PP = append(sd.PP, 0.0)
			}
			sd.Points = append(sd.Points, c.Att8)
			va := s.SubjectVA(subj, nat)
			if va.Error == nil {
				sd.VA = append(sd.VA, va.Pts)
			}
			data[subj] = sd
		}
	}

	keys := []string{}
	for key := range data {
		keys = append(keys, key)
	}
	sort.Sort(sort.StringSlice(keys))

	listData := []subjSummary{}
	for _, key := range keys {
		d := data[key]
		s := subjSummary{Subject: d.Subject, Cohort: d.Cohort}
		s.KS2, err = mean(d.KS2)
		s.HasKS2 = (err == nil)
		pp, _ := mean(d.PP)
		s.PP = 100.0 * pp
		s.Points, _ = mean(d.Points)
		s.VA, err = mean(d.VA)
		s.HasVA = (err == nil)
		s.AvGrade = avGrade(s.Points)
		listData = append(listData, s)
	}

	return listData, nil
}

func mean(data []float64) (float64, error) {

	total := 0.0
	cohort := 0
	for _, d := range data {
		total += d
		cohort++
	}

	if cohort == 0 {
		return 0.0, errors.New("Division by zero")
	}
	return total / float64(cohort), nil
}

// TODO: Change this from hardcoded.
func avGrade(score float64) string {

	grades := []string{"U", "G", "F", "E", "D", "C", "B", "A", "A*"}
	subgrades := []string{"-", "", "+"}

	return grades[int(score+0.5)] + subgrades[int(3*(score+0.5))-3*int(score+0.5)]
}
