package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/andrewcharlton/school-dashboard/analysis/student"
	"github.com/andrewcharlton/school-dashboard/analysis/subject"
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

type pgCellStudent struct {
	UPN  string
	Name string
}

type pgCell struct {
	VA       float64
	Students []pgCellStudent
}

type pgGrid struct {
	Cells    map[string](map[string]pgCell)
	KS2      []string
	Grades   []string
	Counts   map[string]int
	VA       map[string]float64
	TotalVA  float64
	TMExists bool
}

func pgGridAnalysis(subject *subject.Subject, students []student.Student) pgGrid {

	ks2grades := []string{"None", "1", "2", "3C", "3B", "3A", "4C", "4B", "4A", "5C", "5B", "5A", "6"}
	grades := subject.Level.SortedGrades()

	// Initialise grid
	grid := pgGrid{KS2: ks2grades, Grades: grades, Cells: map[string](map[string]pgCell){},
		Counts: map[string]int{}, VA: map[string]float64{}}
	for _, ks2 := range ks2grades {
		grid.Cells[ks2] = map[string]pgCell{}
	}

	// Populate cell lists
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
		cell.Students = append(cell.Students, pgCellStudent{s.UPN, s.Name()})
		grid.Cells[ks2][c.Grd] = cell
		grid.Counts[c.Grd]++
	}

	// Otherwise, assume KS4 and try to load TMs.
	tm, exists := nat.TMs[subject.TM]
	if !exists {
		if subject.TM != "" {
			fmt.Println("TM not found:", subject.TM)
		}
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

type pgStudent struct {
	student.Student
	Class      string
	KS2        string
	Grade      string
	Effort     int
	VAExists   bool
	VA         float64
	Attendance float64
}

func pgStudentList(subject *subject.Subject, students []student.Student) []pgStudent {

	stdnts := []pgStudent{}
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

		stdnts = append(stdnts, pgStudent{
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
