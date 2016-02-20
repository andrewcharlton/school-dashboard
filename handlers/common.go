// Package handlers provides handlers for each of the different
// web pages needed by the dashboard applicaton.
package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/andrewcharlton/school-dashboard/env"
)

// Redirect routes back to the homepage.
func redirect(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		url := "/index/?" + r.URL.RawQuery
		http.Redirect(w, r, url, 301)

	}
}

// Header writes the common html page header and menu bars
func header(e env.Env, w http.ResponseWriter, r *http.Request, detail int) {

	f := getFilter(e, r)
	data := struct {
		School string
		F      database.Filter
		Query  template.URL
	}{
		e.Config.School,
		f,
		template.URL(r.URL.RawQuery),
	}

	err := e.Templates.ExecuteTemplate(w, "header.tmpl", data)
	if err != nil {
		log.Fatal(err)
	}

	filterPage(e, w, r, detail)
}

// Footer writes the common html page header and menu bars
func footer(e env.Env, w http.ResponseWriter, r *http.Request) {

	err := e.Templates.ExecuteTemplate(w, "footer.tmpl", e)
	if err != nil {
		log.Fatal(err)
	}
}

// GetFilter produces a Filter object from the query string
// provided in the http Request
func getFilter(e env.Env, r *http.Request) database.Filter {

	query := r.URL.Query()
	if len(query) == 0 {
		return e.Config.DefaultFilter()
	}

	f := database.Filter{}
	f.Date = query.Get("date")
	f.Resultset = query.Get("resultset")
	f.NatYear = query.Get("natyear")
	f.Year = query.Get("year")
	f.PP = query.Get("pp")
	f.EAL = query.Get("eal")
	f.Gender = query.Get("gender")
	f.SEN = query["sen"]
	f.KS2Bands = query["ks2band"]

	// Change "Other" ethnicity to actual values
	ethnicities := query["ethnicity"]
	other := false
	for _, e := range ethnicities {
		if e == "Other" {
			other = true
		}
		f.Ethnicities = append(f.Ethnicities, e)
	}

	if other {
		for key, oth := range e.OtherEths {
			if oth {
				f.Ethnicities = append(f.Ethnicities, key)
			}
		}
	}
	return f
}

// FilterPage writes the contents of the filter template to w.
// Detail describes the level to which the date should be allowed to be
// filtered:
// 0.  Not at all - National / Resultset / Date only
// 1.  Filter by year
// 2.  Filter by all categories
func filterPage(e env.Env, w http.ResponseWriter, r *http.Request, detail int) {

	f := getFilter(e, r)

	data := struct {
		database.Filter
		Dates       []database.Lookup
		Resultsets  []database.Lookup
		NatYears    []database.Lookup
		Ethnicities []string
		B           map[string]bool //KS2Bands ticked
		E           map[string]bool //Ethnicities checked
		O           map[string]bool //Ethnicities in the "Other" category
		S           map[string]bool //SEN ticked
		Labels      []label
		Detail      int
	}{
		f,
		e.Dates,
		e.Resultsets,
		e.NatYears,
		e.Ethnicities,
		map[string]bool{},
		map[string]bool{},
		e.OtherEths,
		map[string]bool{},
		filterLabels(e, f),
		detail,
	}

	for _, b := range f.KS2Bands {
		data.B[b] = true
	}

	for _, e := range f.Ethnicities {
		data.E[e] = true
	}

	for _, s := range f.SEN {
		data.S[s] = true
	}

	err := e.Templates.ExecuteTemplate(w, "filter.tmpl", data)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}

// Label holds the label text, as well as its formatting options
type label struct {
	Text   string
	Format string
}

// FilterLabels generates the labels for the filter page
func filterLabels(e env.Env, f database.Filter) []label {

	labels := []label{}

	// Lookup date and resultset names
	// Lookup date and resultset names
	date, _ := e.LookupDate(f.Date)
	rs, _ := e.LookupResultset(f.Resultset)
	nat, _ := e.LookupNatYear(f.NatYear)

	labels = append(labels, label{nat, "default"})
	labels = append(labels, label{date, "primary"})
	labels = append(labels, label{rs, "primary"})

	if f.Year != "" {
		labels = append(labels, label{"Yeargroup: " + f.Year, "success"})
	}

	switch f.Gender {
	case "1":
		labels = append(labels, label{"Boys", "warning"})
	case "0":
		labels = append(labels, label{"Girls", "warning"})
	}

	switch f.PP {
	case "1":
		labels = append(labels, label{"Disadvantaged", "warning"})
	case "0":
		labels = append(labels, label{"Non-Disadvantaged", "warning"})
	}

	switch f.EAL {
	case "1":
		labels = append(labels, label{"EAL", "warning"})
	case "0":
		labels = append(labels, label{"Non-EAL", "warning"})
	}

	if len(f.SEN) >= 1 && len(f.SEN) < 4 {
		labels = append(labels, label{"SEN: " + strings.Join(f.SEN, ", "), "warning"})
	}

	if len(f.KS2Bands) >= 1 && len(f.SEN) < 4 {
		labels = append(labels, label{"KS2: " + strings.Join(f.KS2Bands, ", "), "warning"})
	}

	if len(f.Ethnicities) >= 1 && len(f.Ethnicities) <= len(e.Ethnicities) {
		eths := []string{}
		for _, eth := range f.Ethnicities {
			if !e.OtherEths[eth] {
				eths = append(eths, eth)
			}
		}
		labels = append(labels, label{"Ethnicity: " + strings.Join(eths, ", "), "warning"})
	}

	return labels
}

// ChangeYear changes the yeargroup filtered by the query
func changeYear(query url.Values, year string) string {

	if _, exists := query["year"]; exists {
		query.Del("year")
	}
	query.Add("year", year)
	return query.Encode()
}
