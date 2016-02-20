package handlers

import (
	"net/http"
	"strings"

	"github.com/andrewcharlton/school-dashboard/env"
)

// checkRedirect checks whether the query string has
// the correct elements.  If not, it will shorten/lengthen
// as necessary and redirect.
// Detail gives the level of filters that are necessary:
// 0.  None - just date, resultset and natinal year
// 1.  Yeargroup only
// 2.  All filters
func checkRedirect(e env.Env, w http.ResponseWriter, r *http.Request, detail int) bool {

	query := r.URL.Query()
	redirect := false

	// Check that Date, Resultset and NatYear are always present
	for _, key := range []string{"Date", "Resultset", "NatYear"} {
		if _, exists := query[strings.ToLower(key)]; !exists {
			query.Add(strings.ToLower(key), e.Config.Options[key])
			redirect = true
		}
	}

	// If detail is >= 1, ensure a year is present otherwise get rid.
	y, exists := query["year"]
	if detail >= 1 && (!exists || y[0] == "") {
		query.Del("year")
		query.Add("year", e.Config.Options["Year"])
		redirect = true
	}
	if detail == 0 && exists {
		query.Del("year")
		redirect = true
	}

	// If detail <= 1, remove any extra filters
	if detail <= 1 {
		del := []string{}
		for key := range query {
			switch key {
			case "natyear", "date", "resultset", "year":
				continue
			case "name":
				if detail == 1 {
					del = append(del, key)
					redirect = true
				}
			default:
				del = append(del, key)
				redirect = true
			}
		}

		for _, key := range del {
			query.Del(key)
		}
	}

	if !redirect {
		return false
	}

	path := r.URL.Path + "?" + query.Encode()
	http.Redirect(w, r, path, 301)
	return true
}
