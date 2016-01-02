package handlers

import (
	"net/http"

	"github.com/andrewcharlton/school-dashboard/database"
)

// queryOpts details what is needed in the query string
// Year - whether a year group is necessary
// Shory - whether other filter options are needed
type queryOpts struct {
	Year   bool
	Others bool
}

// checkRedirect checks whether the query string has
// the correct elements.  If not, it will shorten/lengthen
// as necessary and redirect.  Returns true if redirect was
// necessary.
func checkRedirect(e database.Env, opts queryOpts, w http.ResponseWriter, r *http.Request) bool {

	query := r.URL.Query()
	redirect := false

	if _, exists := query["date"]; !exists {
		query.Add("date", e.Config.Date)
		redirect = true
	}

	if _, exists := query["resultset"]; !exists {
		query.Add("resultset", e.Config.Resultset)
		redirect = true
	}

	if _, exists := query["natyear"]; !exists {
		query.Add("natyear", e.Config.NatYear)
		redirect = true
	}

	_, exists := query["year"]
	if opts.Year && !exists {
		query.Add("year", e.Config.Year)
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
