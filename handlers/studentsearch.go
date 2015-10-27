package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/andrewcharlton/school-dashboard/env"
)

func SearchRedirect() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		name := query.Get("name")

		query.Del("name")
		url := "/search/" + name + "/?" + query.Encode()

		http.Redirect(w, r, url, 301)
	}
}

func Search(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		Header(e, w, r)
		FilterPage(e, w, r, true)

		name := ""
		path := strings.Split(r.URL.Path, "/")
		for i := len(path) - 1; i > 1; i-- {
			if path[i] != "" {
				name = path[i]
				break
			}
		}

		f := GetFilter(e, r)
		list, err := e.DB.Search(name, f.Date)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		data := struct {
			Name     string
			Query    template.URL
			Results  bool
			Students []database.StudentLookup
		}{
			name,
			template.URL(r.URL.RawQuery),
			(len(list) > 0),
			list,
		}

		err = e.Templates.ExecuteTemplate(w, "studentsearch.html", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}

		Footer(e, w, r)
	}
}
