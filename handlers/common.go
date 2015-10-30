package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/andrewcharlton/school-dashboard/env"
)

// Header writes the common html page header and menu bars
func Header(e env.Env, w http.ResponseWriter, r *http.Request) {

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
func Footer(e env.Env, w http.ResponseWriter, r *http.Request) {

	err := e.Templates.ExecuteTemplate(w, "footer.tmpl", e)
	if err != nil {
		log.Fatal(err)
	}
}

// GetFilter produces a Filter object from the query string
// provided in the http Request
func GetFilter(e env.Env, r *http.Request) database.Filter {

	query := r.URL.Query()
	if len(query) == 0 {
		return e.Config.DefaultFilter()
	}

	f := database.Filter{}
	f.Date = query.Get("date")
	f.Resultset = query.Get("resultset")
	f.Years = query["year"]
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
func FilterPage(e env.Env, w http.ResponseWriter, r *http.Request, short bool) {

	f := GetFilter(e, r)

	data := struct {
		database.Filter
		Dates       []database.Lookup
		Resultsets  []database.Lookup
		Ethnicities []string
		Y           map[string]bool //Years ticked
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
		e.Ethnicities,
		map[string]bool{},
		map[string]bool{},
		map[string]bool{},
		e.OtherEths,
		map[string]bool{},
		FilterLabels(e, f, short),
		short,
	}

	for _, y := range f.Years {
		data.Y[y] = true
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
func FilterLabels(e env.Env, f database.Filter, short bool) []string {

	labels := []string{}

	// Lookup date and resultset names
	// Lookup date and resultset names
	date, rs := "", ""
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

	labels = append(labels, "Date: "+date)
	labels = append(labels, "Resultset: "+rs)

	if short {
		return labels
	}

	if len(f.Years) >= 1 {
		labels = append(labels, "Years: "+strings.Join(f.Years, ", "))
	}

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

// PageNotFound writes a 'Page not found error'
func PageNotFound(w http.ResponseWriter) {

	fmt.Fprintf(w, "<h4>Error: 404</h4><br>Page not found.")

}
