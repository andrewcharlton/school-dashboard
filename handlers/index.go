package handlers

import (
	"fmt"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/andrewcharlton/school-dashboard/env"
)

// Index produces a landing page
func Index(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		header(e, w, r, 0)
		defer footer(e, w, r)

		news, err := e.News()
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}

		data := struct {
			School string
			News   []database.NewsItem
		}{
			e.Config.School,
			news,
		}

		err = e.Templates.ExecuteTemplate(w, "index.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}

	}
}
