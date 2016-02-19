package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/andrewcharlton/school-dashboard/analysis/group"
	"github.com/andrewcharlton/school-dashboard/analysis/subject"
	"github.com/andrewcharlton/school-dashboard/env"
)

func SubjectGroups(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := strings.Split(r.URL.Path, "/")
		switch len(path) {
		case 3:
			selectSubject(e, w, r, "Progress Grid")
		case 4:
			selectLevel(e, w, r, "Progress Grid")
		case 5:
			selectYear(e, w, r, "Progress Grid")
		case 6:
			progressGridPage(e, w, r)
		}
	}
}

func subjectGroupPage(e env.Env, w http.ResponseWriter, r *http.Request) {

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
	subj := e.Subjects[subjID]

	f := getFilter(e, r)
	g, err := e.GroupByFilteredClass(path[3], "", f)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}

	classes, err := e.Classes(path[3], f.Date)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
	sort.Sort(sort.StringSlice(classes))

	clsGrps := []subGroup{}
	for _, c := range classes {
		clsGrps = append(clsGrps, subGroup{c, template.URL(c), g.SubGroup(group.Class(subj.Subj, c))})
	}

	data := struct {
		Subj      *subject.Subject
		SubGroups []subGroup
		Classes   []subGroup
		Query     template.URL
	}{
		subj,
		subGroups(g),
		clsGrps,
		template.URL(r.URL.RawQuery),
	}

	err = e.Templates.ExecuteTemplate(w, "subjectgroups.tmpl", data)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

}
