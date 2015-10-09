package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/andrewcharlton/school-dashboard/database"
)

// Option holds the details of a table lookup which contains
// an id and a name
type Option struct {
	ID   string
	Name string
}

// GetFilter returns a Filter struct from the query options
// on a page.
func (e Env) GetFilter(query url.Values) (database.Filter, error) {

	if len(query) == 0 {
		f, err := e.DB.DefaultFilter()
		if err != nil {
			return database.Filter{}, err
		}
		return f, nil
	}

	f := database.Filter{}
	f.Date = query.Get("date")
	f.Resultset = query.Get("resultset")
	f.Exams = query.Get("exams")
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

	return f, nil

}

// FilterPage writes the contents of the filter template to w.
func (e Env) FilterPage(w http.ResponseWriter, f database.Filter, short bool) error {

	data := struct {
		database.Filter
		Dates       []Option
		Resultsets  []Option
		Ethnicities []string
		Y           map[string]bool //Years ticked
		B           map[string]bool //KS2Bands ticked
		E           map[string]bool //Ethnicities checked
		O           map[string]bool //Ethnicities in the "Other" category
		S           map[string]bool //SEN ticked
		Labels      []string
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
		e.FilterLabels(f, short),
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

	var t string
	switch {
	case short:
		t = "filter-short.tmpl"
	case !short:
		t = "filter.tmpl"
	}

	err := e.Templates.ExecuteTemplate(w, t, data)
	if err != nil {
		return err
	}

	return nil
}

// FilterLabels generates the labels for the filter page
func (e Env) FilterLabels(f database.Filter, short bool) []string {

	labels := []string{}

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
	labels = append(labels, "Exams: "+f.Exams)

	if short {
		return labels
	}

	if len(f.Years) >= 1 {
		labels = append(labels, "Years: "+strings.Join(f.Years, ", "))
	}

	switch {
	case f.Gender == "1":
		labels = append(labels, "Boys")
	case f.Gender == "0":
		labels = append(labels, "Girls")
	}

	switch {
	case f.PP == "1":
		labels = append(labels, "Disadvanted")
	case f.PP == "0":
		labels = append(labels, "Non-Disadvantaged")
	}

	switch {
	case f.EAL == "1":
		labels = append(labels, "EAL")
	case f.EAL == "0":
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
