package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"

	"github.com/andrewcharlton/school-dashboard/analysis/group"
	"github.com/andrewcharlton/school-dashboard/analysis/subject"
	"github.com/andrewcharlton/school-dashboard/env"
)

// KS3Summary produces a page with the student
func KS3Summary(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if redir := checkRedirect(e, w, r, 2); redir {
			return
		}
		header(e, w, r, 2)
		defer footer(e, w, r)

		f := getFilter(e, r)
		g, err := e.GroupByFilter(f)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		ks3Subjects := subject.SubjectList{}
		var lvl subject.Level
		for _, s := range e.Subjects {
			if s.Lvl == "KS3" {
				lvl = s.Level
				ks3Subjects = append(ks3Subjects, *s)
			}
		}
		sort.Sort(ks3Subjects)

		data := struct {
			Query    template.URL
			Subjects subject.SubjectList
			KS3      subject.Level
			Group    group.Group
		}{
			template.URL(r.URL.RawQuery),
			ks3Subjects,
			lvl,
			g,
		}

		err = e.Templates.ExecuteTemplate(w, "ks3summary.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}
