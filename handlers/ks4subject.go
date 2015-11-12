package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/andrewcharlton/school-dashboard/analysis"
	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/andrewcharlton/school-dashboard/national"
)

func KS4Subject(e database.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		Header(e, w, r)
		FilterPage(e, w, r, false)
		defer Footer(e, w, r)

		path := strings.Split(r.URL.Path, "/")
		switch {
		case len(path) < 3 || path[2] == "":
			ks4Subjects(e, w, r)
		case len(path) < 4 || path[3] == "":
			ks4Classes(e, w, r, path[2])
		default:
			ks4Analysis(e, w, r, path[2], path[3])
		}
	}
}

// Returns page of subjects to pick from.
func ks4Subjects(e database.Env, w http.ResponseWriter, r *http.Request) {

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

	err = e.Templates.ExecuteTemplate(w, "ks4subject-subjects.tmpl", data)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}

// Returns page of classes in that subject to pick from.
func ks4Classes(e database.Env, w http.ResponseWriter, r *http.Request, subj string) {

	f := GetFilter(e, r)
	classes, err := e.DB.Classes(subj, f.Date)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	y10, y11 := []string{"All Students"}, []string{"All Students"}
	for _, class := range classes {
		if strings.HasPrefix(class, "10") {
			y10 = append(y10, class)
		}
		if strings.HasPrefix(class, "11") {
			y11 = append(y11, class)
		}
	}

	data := struct {
		Subject  string
		Y10      []string
		Y11      []string
		Query    template.URL
		Y10Query template.URL
		Y11Query template.URL
	}{
		subj,
		y10,
		y11,
		template.URL(r.URL.RawQuery),
		template.URL(ChangeYear(r.URL.Query(), "10")),
		template.URL(ChangeYear(r.URL.Query(), "11")),
	}

	err = e.Templates.ExecuteTemplate(w, "ks4subject-classes.tmpl", data)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}

// Performs analysis of the results
func ks4Analysis(e database.Env, w http.ResponseWriter, r *http.Request, subj, class string) {

	f := GetFilter(e, r)
	cls := class
	if class == "All Students" {
		cls = ""
	}

	g, err := e.DB.GroupByFilteredClass(subj, cls, f)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}

	nat := e.Nationals[f.NatYear]
	type student struct {
		UPN        string
		Name       string
		Class      string
		KS2        string
		Grade      string
		Effort     int
		VA         float64
		VAExists   bool
		Attendance float64
	}

	students := []student{}
	groupVA := float64(0)
	num := 0
	for _, s := range g.Students {
		stdnt := student{UPN: s.UPN, Name: s.Name(), Attendance: s.Attendance.Latest()}
		c, exists := s.Courses[subj]
		if !exists {
			continue
		}
		stdnt.Class = c.Class
		stdnt.Grade = c.Grd
		stdnt.Effort = c.Effort
		switch c.KS2Prior {
		case "En":
			stdnt.KS2 = s.KS2.En
		case "Ma":
			stdnt.KS2 = s.KS2.Ma
		default:
			stdnt.KS2 = s.KS2.Av
		}
		va := s.SubjectVA(subj, nat)
		if va.Error == nil {
			stdnt.VAExists = true
			stdnt.VA = va.Pts
			groupVA += va.Pts
			num += 1
		}
		students = append(students, stdnt)
	}

	data := struct {
		Subject  string
		Class    string
		Students []student
		Group    analysis.Group
		Query    template.URL
		Nat      national.National
	}{
		subj,
		class,
		students,
		g,
		template.URL(r.URL.RawQuery),
		nat,
	}

	err = e.Templates.ExecuteTemplate(w, "ks4subject-analysis.tmpl", data)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}
