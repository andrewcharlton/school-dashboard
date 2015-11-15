package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/andrewcharlton/school-dashboard/analysis"
	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/andrewcharlton/school-dashboard/national"
)

func KS4Subject(e database.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := strings.Split(r.URL.Path, "/")
		switch len(path) {
		case 3:
			selectSubject(e, w, r, "KS4 Analysis")
		case 4:
			selectLevel(e, w, r, "KS4 Analysis")
		case 5:
			selectClass(e, w, r, "KS4 Analysis")
		case 6:
			ks4SubjectAnalysis(e, w, r)
		}
	}
}

// Performs analysis of the results
func ks4SubjectAnalysis(e database.Env, w http.ResponseWriter, r *http.Request) {

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
	if path[4] == "All Students" {
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
		Students []vaStudent
		Grid     vaGrid
		Query    template.URL
	}{
		subject.Subj,
		subject.Lvl,
		path[3],
		path[4],
		ks4VAStudents(subject, g.Students, nat),
		ks4VAGrid(subject, g.Students, nat),
		template.URL(ShortenQuery(e, r.URL.Query())),
	}

	err = e.Templates.ExecuteTemplate(w, "ks4subject.tmpl", data)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

}

type vaCellStudent struct {
	UPN  string
	Name string
}

type vaCell struct {
	VA       float64
	Students []vaCellStudent
}

type vaGrid struct {
	Cells    map[string](map[string]vaCell)
	KS2      []string
	Grades   []string
	VA       map[string]float64
	TotalVA  float64
	TMExists bool
}

func ks4VAGrid(subject *analysis.Subject, students []analysis.Student, nat national.National) vaGrid {

	ks2grades := []string{"None", "1", "2", "3C", "3B", "3A", "4C", "4B", "4A", "5C", "5B", "5A", "6"}
	grades := subject.Level.SortedGrades()

	grid := vaGrid{KS2: ks2grades, Grades: grades, Cells: map[string](map[string]vaCell){},
		VA: map[string]float64{}}
	for _, ks2 := range ks2grades {
		grid.Cells[ks2] = map[string]vaCell{}
	}

	for _, s := range students {
		c, exists := s.Courses[subject.Subj]
		if !exists {
			continue
		}

		var ks2 string
		switch subject.KS2Prior {
		case "En":
			ks2 = s.KS2.En
		case "Ma":
			ks2 = s.KS2.Ma
		default:
			ks2 = s.KS2.Av
		}

		if ks2 == "" {
			ks2 = "None"
		}

		cell := grid.Cells[ks2][c.Grd]
		cell.Students = append(cell.Students, vaCellStudent{s.UPN, s.Name()})
		grid.Cells[ks2][c.Grd] = cell
	}

	tm, exists := nat.TMs[subject.TM]
	if !exists {
		return grid
	}
	grid.TMExists = true

	// Calculate VAs
	totalVA := 0.0
	totalN := 0
	for _, ks2 := range ks2grades {
		if ks2 == "None" {
			continue
		}
		rowVA := 0.0
		rowN := 0
		for _, grd := range grades {
			va, err := tm.ValueAdded(ks2, grd)
			if ks2 == "None" {
				continue
			}
			if err != nil {
				fmt.Println("KS4VAGrid - Error:", err)
			}
			cell := grid.Cells[ks2][grd]
			cell.VA = va
			grid.Cells[ks2][grd] = cell
			rowVA += va * float64(len(grid.Cells[ks2][grd].Students))
			rowN += len(grid.Cells[ks2][grd].Students)
		}
		totalVA += rowVA
		totalN += rowN
		if rowN != 0 {
			grid.VA[ks2] = rowVA / float64(rowN)
		}
	}
	grid.TotalVA = totalVA / float64(totalN)

	return grid
}

type vaStudent struct {
	analysis.Student
	Class      string
	KS2        string
	Grade      string
	Effort     int
	VAExists   bool
	VA         float64
	Attendance float64
}

func ks4VAStudents(subject *analysis.Subject, students []analysis.Student, nat national.National) []vaStudent {

	stdnts := []vaStudent{}
	for _, s := range students {
		c, exists := s.Courses[subject.Subj]
		if !exists {
			continue
		}

		ks2 := ""
		switch subject.KS2Prior {
		case "En":
			ks2 = s.KS2.En
		case "Ma":
			ks2 = s.KS2.Ma
		default:
			ks2 = s.KS2.Av
		}

		va := s.SubjectVA(subject.Subj, nat)

		stdnts = append(stdnts, vaStudent{
			s,
			c.Class,
			ks2,
			c.Grd,
			c.Effort,
			va.Error == nil,
			va.Pts,
			s.Attendance.Latest(),
		})
	}

	return stdnts
}
