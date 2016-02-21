package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/andrewcharlton/school-dashboard/analysis/group"
	"github.com/andrewcharlton/school-dashboard/env"
)

// ProgressGrid produces a page containing a progress grid for a subject
func ProgressGrid(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := strings.Split(r.URL.Path, "/")
		switch len(path) {
		case 3:
			selectSubject(e, w, r, "Progress Grid")
		case 4:
			selectLevel(e, w, r, "Progress Grid")
		case 5:
			selectClass(e, w, r, "Progress Grid")
		case 6:
			progressGridPage(e, w, r)
		}
	}
}

// Performs analysis of the results
func progressGridPage(e env.Env, w http.ResponseWriter, r *http.Request) {

	if redir := checkRedirect(e, w, r, 2); redir {
		return
	}
	header(e, w, r, 2)
	defer footer(e, w, r)

	path := strings.Split(r.URL.Path, "/")
	subjID, err := strconv.Atoi(path[3])
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	subject := e.Subjects[subjID]
	class := path[4]
	if strings.HasPrefix(path[4], "All") {
		class = ""
	}

	f := getFilter(e, r)
	g, err := e.GroupByFilteredClass(path[3], class, f)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}

	data := struct {
		Query        template.URL
		Year         string
		Subject      string
		Level        string
		SubjID       string
		Class        string
		KS2Prior     string
		Group        group.Group
		ProgressGrid group.ProgressGrid
	}{
		template.URL(r.URL.RawQuery),
		f.Year,
		subject.Subj,
		subject.Lvl,
		path[3],
		path[4],
		subject.KS2Prior,
		g,
		g.ProgressGrid(subject, f.NatYear),
	}

	err = e.Templates.ExecuteTemplate(w, "progressgrid.tmpl", data)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

}
