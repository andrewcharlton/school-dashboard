package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/andrewcharlton/school-dashboard/env"
)

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
			pgAnalysis(e, w, r)
		}
	}
}

// Performs analysis of the results
func pgAnalysis(e env.Env, w http.ResponseWriter, r *http.Request) {

	if redir := checkRedirect(e, queryOpts{true, true}, w, r); redir {
		return
	}

	Header(e, w, r)
	FilterPage(e, w, r, false)
	defer Footer(e, w, r)

	path := strings.Split(r.URL.Path, "/")
	subjID, err := strconv.Atoi(path[3])
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	subject := e.DB.Subjects()[subjID]
	class := path[4]
	if strings.HasPrefix(path[4], "All") {
		class = ""
	}

	f := GetFilter(e, r)
	g, err := e.DB.GroupByFilteredClass(path[3], class, f)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}

	nat := e.Nationals[f.NatYear]

	data := struct {
		Subject  string
		Level    string
		SubjID   string
		Class    string
		Students []pgStudent
		Grid     pgGrid
		Query    template.URL
	}{
		subject.Subj,
		subject.Lvl,
		path[3],
		path[4],
		pgStudentList(subject, g.Students, nat),
		pgGridAnalysis(subject, g.Students, nat),
		template.URL(r.URL.RawQuery),
	}

	err = e.Templates.ExecuteTemplate(w, "progressgrid.tmpl", data)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

}
