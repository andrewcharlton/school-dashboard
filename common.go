package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/andrewcharlton/school-dashboard/database"
)

// Holds details of items in menus
type menuItem struct {
	Sep  bool
	Name string
	URL  template.URL
}

// Header writes the contents of the header template to w.
func (e Env) Header(w http.ResponseWriter, r *http.Request) error {

	f, err := e.GetFilter(r.URL.Query())
	if err != nil {
		return err
	}
	q := strings.Split(r.URL.RawQuery, "&search")[0]

	data := struct {
		database.Filter
		Menu      []menuItem
		Dashboard template.URL
	}{
		f,
		[]menuItem{},
		template.URL("../dashboard/?" + q),
	}

	links := []string{"Headlines", "", "Data"}
	for _, l := range links {
		if l == "" {
			data.Menu = append(data.Menu, menuItem{Sep: true})
		} else {
			data.Menu = append(data.Menu, menuItem{false, l, template.URL("../" + l + "/?" + q)})
		}
	}

	err = e.Templates.ExecuteTemplate(w, "header.tmpl", data)
	if err != nil {
		return err
	}
	return nil
}

// Footer writes the contents of the footer template to w.
func (e Env) Footer(w http.ResponseWriter) error {

	err := e.Templates.ExecuteTemplate(w, "footer.tmpl", nil)
	if err != nil {
		return err
	}
	return nil
}
