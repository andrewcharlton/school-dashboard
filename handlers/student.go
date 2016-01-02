package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/andrewcharlton/school-dashboard/analysis"
	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/andrewcharlton/school-dashboard/national"
)

func Student(e database.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if redir := checkRedirect(e, queryOpts{false, false}, w, r); redir {
			return
		}

		Header(e, w, r)
		FilterPage(e, w, r, true)
		defer Footer(e, w, r)

		upn := ""
		path := strings.Split(r.URL.Path, "/")
		for i := len(path) - 1; i > 1; i-- {
			if path[i] != "" {
				upn = path[i]
				break
			}
		}

		query := r.URL.Query()
		f := database.StudentFilter{
			upn,
			query.Get("date"),
			query.Get("resultset"),
		}

		s, err := e.DB.Student(f)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		nat := e.Nationals[query.Get("natyear")]
		p8, err := nat.Progress8(s.KS2.APS)

		data := struct {
			Student analysis.Student
			Nat     national.Progress8
			HasNat  bool
		}{
			s,
			p8,
			err == nil,
		}

		err = e.Templates.ExecuteTemplate(w, "student.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}
