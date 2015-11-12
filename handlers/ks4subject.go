package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/andrewcharlton/school-dashboard/analysis"
	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/andrewcharlton/school-dashboard/level"
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

	var tm national.TransitionMatrix
	for _, s := range g.Students {
		t, err := s.TM(subj, nat)
		if err == nil {
			tm = t
			break
		}
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

	grid, va := ks4Grid(g, subj, tm)
	data := struct {
		Subject  string
		Class    string
		Students []student
		Group    analysis.Group
		Grid     progressGrid
		VA       float64
		Query    template.URL
		Nat      national.National
	}{
		subj,
		class,
		students,
		g,
		grid,
		va,
		template.URL(r.URL.RawQuery),
		nat,
	}

	err = e.Templates.ExecuteTemplate(w, "ks4subject-analysis.tmpl", data)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}

type cell struct {
	Num     int
	BGColor string
	FGColor string
}

type progressGridRow struct {
	KS2   string
	Cells []cell
	VA    float64
}

type progressGrid struct {
	Rows   []progressGridRow
	Grades []string
}

func ks4Grid(g analysis.Group, subj string, tm national.TransitionMatrix) (progressGrid, float64) {

	var lvl *level.Level
	nums := map[string](map[string]int){}
	for _, s := range g.Students {
		c, exists := s.Courses[subj]
		if !exists {
			continue
		}

		if lvl == nil {
			lvl = c.Level
		}

		var ks2 string
		switch c.KS2Prior {
		case "En":
			ks2 = s.KS2.En
		case "Ma":
			ks2 = s.KS2.Ma
		default:
			ks2 = s.KS2.Av
		}

		row, exists := nums[ks2]
		if !exists {
			row = map[string]int{}
		}
		row[c.Grd] += 1
		nums[ks2] = row
	}

	totalVA := 0.0
	totalN := 0

	levelGrades := lvl.SortedGrades()
	grid := progressGrid{Grades: levelGrades, Rows: []progressGridRow{}}
	ks2Grades := []string{"1", "2", "3C", "3B", "3A", "4C", "4B", "4A", "5C", "5B", "5A", "6"}

	for _, ks2 := range ks2Grades {
		row := progressGridRow{KS2: ks2, Cells: []cell{}}
		n := 0
		for _, grade := range levelGrades {
			c := cell{Num: nums[ks2][grade], FGColor: "#000000", BGColor: "#FFFFFF"}
			va, err := tm.ValueAdded(ks2, grade)
			if err == nil {
				switch {
				case va < -0.33:
					c.BGColor = "#FA6E2D"
				case va > 0.0:
					c.BGColor = "#2EB02E"
				default:
					c.BGColor = "#FCF4A4"
				}
				row.VA += va * float64(c.Num)
				totalVA += va * float64(c.Num)
				n += c.Num
				totalN += 1
			}
			row.Cells = append(row.Cells, c)
		}
		if n > 0 {
			row.VA = row.VA / float64(n)
		}
		grid.Rows = append(grid.Rows, row)
	}

	if totalN == 0 {
		return grid, 0.0
	}
	return grid, totalVA / float64(totalN)
}
