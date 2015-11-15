package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/andrewcharlton/school-dashboard/analysis"
	"github.com/andrewcharlton/school-dashboard/database"
)

func ClassList(e database.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := strings.Split(r.URL.Path, "/")
		switch len(path) {
		case 3:
			selectSubject(e, w, r, "Class Lists")
		case 4:
			selectLevel(e, w, r, "Class Lists")
		case 5:
			selectClass(e, w, r, "Class Lists")
		case 6:
			classStudentList(e, w, r)
		}
	}
}

// Class list for the students
func classStudentList(e database.Env, w http.ResponseWriter, r *http.Request) {

	Header(e, w, r)
	FilterPage(e, w, r, true)
	defer Footer(e, w, r)

	path := strings.Split(r.URL.Path, "/")
	subjID, err := strconv.Atoi(path[3])
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	subject := e.DB.Subjects()[subjID]
	class := path[4]

	f := GetFilter(e, r)
	g, err := e.DB.GroupByClass(path[3], class, f)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	data := struct {
		Subject  string
		Level    string
		SubjID   string
		Class    string
		Query    template.URL
		Students []analysis.Student
	}{
		subject.Subj,
		subject.Lvl,
		path[3],
		class,
		template.URL(ShortenQuery(e, r.URL.Query())),
		g.Students,
	}

	err = e.Templates.ExecuteTemplate(w, "classlist.tmpl", data)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}
