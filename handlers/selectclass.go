package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/andrewcharlton/school-dashboard/env"
)

// These functions provide pages to select a subject/class for analysis etc.
// Expects URL path in the form:
// /basepath/*subject name*/*subject id*/*class name*/?*query*

// Produce page to pick a subject from
func selectSubject(e env.Env, w http.ResponseWriter, r *http.Request, heading string) {

	if redir := checkRedirect(e, queryOpts{false, false}, w, r); redir {
		return
	}

	Header(e, w, r)
	FilterPage(e, w, r, true)
	defer Footer(e, w, r)

	subjects := e.DB.Subjects()

	distinct := map[string]struct{}{}
	for _, subj := range subjects {
		distinct[subj.Subj] = struct{}{}
	}

	data := struct {
		Heading  string
		Subjects []string
		Path     template.URL
		Query    template.URL
	}{
		heading,
		[]string{},
		template.URL(r.URL.Path),
		template.URL(r.URL.RawQuery),
	}

	for subj := range distinct {
		data.Subjects = append(data.Subjects, subj)
	}
	sort.Sort(sort.StringSlice(data.Subjects))

	err := e.Templates.ExecuteTemplate(w, "select-subject.tmpl", data)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}

type subjLevel struct {
	SubjID int
	Level  string
}

type subjLevels []subjLevel

func (s subjLevels) Len() int           { return len(s) }
func (s subjLevels) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s subjLevels) Less(i, j int) bool { return s[i].Level < s[j].Level }

// Produce page to pick a level from
func selectLevel(e env.Env, w http.ResponseWriter, r *http.Request, heading string) {

	if redir := checkRedirect(e, queryOpts{false, false}, w, r); redir {
		return
	}

	// Assume subject is final part of path
	path := strings.Split(r.URL.Path, "/")
	subject := path[2]

	levels := subjLevels{}
	for id, s := range e.DB.Subjects() {
		if s.Subj == subject {
			levels = append(levels, subjLevel{id, s.Lvl})
		}
	}
	sort.Sort(levels)

	// If only one level, just redirect to that page
	if len(levels) == 1 {
		url := fmt.Sprintf("%v/%v/?%v", r.URL.Path, levels[0].SubjID, r.URL.RawQuery)
		http.Redirect(w, r, url, 301)
		return
	}

	Header(e, w, r)
	FilterPage(e, w, r, true)
	defer Footer(e, w, r)

	data := struct {
		Heading  string
		Subject  string
		Levels   subjLevels
		BasePath template.URL
		Path     template.URL
		Query    template.URL
	}{
		heading,
		subject,
		levels,
		template.URL("/" + path[1]),
		template.URL(r.URL.Path),
		template.URL(r.URL.RawQuery),
	}

	err := e.Templates.ExecuteTemplate(w, "select-level.tmpl", data)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}

// Produce page to pick a class from
func selectClass(e env.Env, w http.ResponseWriter, r *http.Request, heading string) {

	if redir := checkRedirect(e, queryOpts{false, false}, w, r); redir {
		return
	}

	Header(e, w, r)
	FilterPage(e, w, r, true)
	defer Footer(e, w, r)

	// Assume subject name and subj_id are last two parts of the path
	path := strings.Split(r.URL.Path, "/")
	subject := path[2]
	subjID, err := strconv.Atoi(path[3])
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	level := e.DB.Subjects()[subjID].Lvl

	f := GetFilter(e, r)
	classes, err := e.DB.Classes(path[3], f.Date)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}

	data := struct {
		Heading  string
		Subject  string
		Level    string
		Years    []string
		Classes  map[string]([]string)
		Queries  map[string]template.URL
		BasePath template.URL
		Path     template.URL
		Query    template.URL
	}{
		heading,
		subject,
		level,
		[]string{},
		map[string]([]string){},
		map[string]template.URL{},
		template.URL("/" + path[1]),
		template.URL(r.URL.Path),
		template.URL(r.URL.RawQuery),
	}

	years := map[string]bool{}
	for _, class := range classes {
		for _, year := range []string{"7", "8", "9", "10", "11"} {
			if strings.HasPrefix(class, year) {
				years[year] = true
				data.Classes[year] = append(data.Classes[year], class)
			}
		}
	}

	for _, year := range []string{"7", "8", "9", "10", "11"} {
		if years[year] {
			data.Years = append(data.Years, year)
			data.Queries[year] = template.URL(ChangeYear(r.URL.Query(), year))
		}
	}

	err = e.Templates.ExecuteTemplate(w, "select-class.tmpl", data)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}

}
