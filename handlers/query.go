package handlers

import (
	"net/http"
	"strings"

	"github.com/andrewcharlton/school-dashboard/env"
)

// queryOpts details what is needed in the query string
// Year - whether a year group is necessary
// Others - whether other filter options are needed
type queryOpts struct {
	Year   bool
	Others bool
}

// checkRedirect checks whether the query string has
// the correct elements.  If not, it will shorten/lengthen
// as necessary and redirect.  Returns true if redirect was
// necessary.
func checkRedirect(e env.Env, opts queryOpts, w http.ResponseWriter, r *http.Request) bool {

	query := r.URL.Query()
	redirect := false

	for _, key := range []string{"Date", "Resultset", "NatYear"} {
		if _, exists := query[strings.ToLower(key)]; !exists {
			query.Add(strings.ToLower(key), e.Config.Options[key])
			redirect = true
		}
	}

	y, exists := query["year"]
	if opts.Year && (!exists || y[0] == "") {
		query.Del("year")
		query.Add("year", e.Config.Options["Year"])
		redirect = true
	}
	if !opts.Year && exists {
		query.Del("year")
		redirect = true
	}

	if !opts.Others {
		del := []string{}
		for key, _ := range query {
			switch key {
			case "natyear", "date", "resultset", "year":
				continue
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
