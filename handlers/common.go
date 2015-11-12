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
)

// Redirect routes back to the homepage.
func Redirect(e database.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Redirect being called!")
		query := ShortenQuery(e, r.URL.Query())
		url := "/index/?" + query
		http.Redirect(w, r, url, 301)

	}
}

// Header writes the common html page header and menu bars
func Header(e database.Env, w http.ResponseWriter, r *http.Request) {

	f := GetFilter(e, r)
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

}

// Footer writes the common html page header and menu bars
func Footer(e database.Env, w http.ResponseWriter, r *http.Request) {

	err := e.Templates.ExecuteTemplate(w, "footer.tmpl", e)
	if err != nil {
		log.Fatal(err)
	}
}

// GetFilter produces a Filter object from the query string
// provided in the http Request
func GetFilter(e database.Env, r *http.Request) database.Filter {

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
func FilterPage(e database.Env, w http.ResponseWriter, r *http.Request, short bool) {

	f := GetFilter(e, r)

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
		Labels      []string
		Short       bool
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
		FilterLabels(e, f, short),
		short,
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

// FilterLabels generates the labels for the filter page
func FilterLabels(e database.Env, f database.Filter, short bool) []string {

	labels := []string{}

	// Lookup date and resultset names
	// Lookup date and resultset names
	var date, rs, nat string
	for _, d := range e.Dates {
		if d.ID == f.Date {
			date = d.Name
			break
		}
	}
	for _, r := range e.Resultsets {
		if r.ID == f.Resultset {
			rs = r.Name
		}
	}
	for _, r := range e.NatYears {
		if r.ID == f.NatYear {
			nat = r.Name
		}
	}

	labels = append(labels, "Date: "+date)
	labels = append(labels, "Resultset: "+rs)
	labels = append(labels, "National: "+nat)

	if short {
		return labels
	}

	labels = append(labels, "Yeargroup: "+f.Year)

	switch f.Gender {
	case "1":
		labels = append(labels, "Boys")
	case "0":
		labels = append(labels, "Girls")
	}

	switch f.PP {
	case "1":
		labels = append(labels, "Disadvantaged")
	case "0":
		labels = append(labels, "Non-Disadvantaged")
	}

	switch f.EAL {
	case "1":
		labels = append(labels, "EAL")
	case "0":
		labels = append(labels, "Non-EAL")
	}

	if len(f.SEN) >= 1 && len(f.SEN) < 4 {
		labels = append(labels, "SEN: "+strings.Join(f.SEN, ", "))
	}

	if len(f.KS2Bands) >= 1 && len(f.SEN) < 4 {
		labels = append(labels, "KS2: "+strings.Join(f.KS2Bands, ", "))
	}

	if len(f.Ethnicities) >= 1 && len(f.Ethnicities) <= len(e.Ethnicities) {
		eths := []string{}
		for _, eth := range f.Ethnicities {
			if !e.OtherEths[eth] {
				eths = append(eths, eth)
			}
		}
		labels = append(labels, "Ethnicity: "+strings.Join(eths, ", "))
	}

	return labels
}

// ShortenQuery removes all group filtering options, leaving
// just date, resultset and national/yeargroup options.
// If these are not present, adds in default options.
func ShortenQuery(e database.Env, query url.Values) string {

	del := []string{}
	for key, _ := range query {
		switch key {
		case "date", "resultset", "natyear", "year":
			continue
		default:
			del = append(del, key)
		}
	}

	for _, key := range del {
		query.Del(key)
	}

	return LengthenQuery(e, query)
}

// LengthenQuery adds any default filtering options if they are
// not already present.
func LengthenQuery(e database.Env, query url.Values) string {

	if _, exists := query["date"]; !exists {
		query.Add("date", e.Config.Date)
	}

	if _, exists := query["resultset"]; !exists {
		query.Add("resultset", e.Config.Resultset)
	}

	if _, exists := query["natyear"]; !exists {
		query.Add("natyear", e.Config.NatYear)
	}

	if _, exists := query["year"]; !exists {
		query.Add("year", e.Config.Year)
	}

	return query.Encode()
}

// ChangeYear changes the yeargroup filtered by the query
func ChangeYear(query url.Values, year string) string {

	if _, exists := query["year"]; exists {
		query.Del("year")
	}
	query.Add("year", year)
	return query.Encode()
}
