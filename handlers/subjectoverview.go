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

type subjSummary struct {
	Subject *subject.Subject
	Group   group.Group
}

type subjSummaries []subjSummary

func (s subjSummaries) Len() int           { return len(s) }
func (s subjSummaries) Less(i, j int) bool { return s[i].Subject.Subj < s[j].Subject.Subj }
func (s subjSummaries) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// SubjectOverview provides a list of subjects with
func SubjectOverview(e env.Env) http.HandlerFunc {
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

		summaries := subjSummaries{}
		for _, subj := range e.Subjects {
			subGroup := g.SubGroup(group.Studying(subj.Subj, subj.SubjID))
			if len(subGroup.Students) == 0 {
				continue
			}
			summaries = append(summaries, subjSummary{subj, subGroup})
		}
		sort.Sort(summaries)

		data := struct {
			Query     template.URL
			Year      string
			Summaries subjSummaries
		}{
			template.URL(r.URL.RawQuery),
			f.Year,
			summaries,
		}

		err = e.Templates.ExecuteTemplate(w, "subject-overview.tmpl", data)
		if err != nil {
			fmt.Println(err)
		}
	}
}
