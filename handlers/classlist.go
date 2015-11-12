package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/andrewcharlton/school-dashboard/analysis"
	"github.com/andrewcharlton/school-dashboard/database"
)

func ClassList(e database.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		Header(e, w, r)
		FilterPage(e, w, r, true)
		defer Footer(e, w, r)

		path := strings.Split(r.URL.Path, "/")

		switch {
		case len(path) < 3 || path[2] == "":
			clSubjects(e, w, r)
		case len(path) < 4 || path[3] == "":
			clClasses(e, w, r, path[2])
		default:
			clStudents(e, w, r, path[2], path[3])
		}
	}
}

// Return page to pick a subject from.
func clSubjects(e database.Env, w http.ResponseWriter, r *http.Request) {

	subjects, err := e.DB.Subjects()
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	data := struct {
		Subjects []string
		Query    template.URL
	}{
		subjects,
		template.URL(r.URL.RawQuery),
	}

	err = e.Templates.ExecuteTemplate(w, "classlist-subjects.tmpl", data)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}

// Return page to pick a class from.
func clClasses(e database.Env, w http.ResponseWriter, r *http.Request, subj string) {

	f := GetFilter(e, r)
	classes, err := e.DB.Classes(subj, f.Date)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	data := struct {
		Subject string
		Classes []string
		Query   template.URL
	}{
		subj,
		classes,
		template.URL(r.URL.RawQuery),
	}

	err = e.Templates.ExecuteTemplate(w, "classlist-classes.tmpl", data)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}

// Return a list of students
func clStudents(e database.Env, w http.ResponseWriter, r *http.Request, subj, class string) {

	f := GetFilter(e, r)
	g, err := e.DB.GroupByClass(subj, class, f)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	data := struct {
		Subject  string
		Class    string
		Query    template.URL
		Students []analysis.Student
	}{
		subj,
		class,
		template.URL(ShortenQuery(e, r.URL.Query())),
		g.Students,
	}

	err = e.Templates.ExecuteTemplate(w, "classlist-students.tmpl", data)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}
