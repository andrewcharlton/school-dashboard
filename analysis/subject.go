package analysis

import (
	"fmt"

	"github.com/andrewcharlton/school-dashboard/national"
)

// A PGCellStudent holds the name and UPN of a student, used for hovering
// over cells.
type PGCellStudent struct {
	UPN  string
	Name string
}

// A PGCell holds the VA for that combination of KS2 and achieved grade,
// and the details of the students at that grade.
type PGCell struct {
	VA       float64
	Students []PGCellStudent
}

// A ProgressGrid showing the grades achieved by each student, broken down by
// prior attainment at KS2.  VA scores are also held for each row.
type ProgressGrid struct {
	// Cells are referenced by KS2, then by Grade achieved.
	Cells    map[string](map[string]PGCell)
	KS2      []string
	Grades   []string
	Counts   map[string]int
	VA       map[string]float64
	TotalVA  float64
	TMExists bool
}

// PGAnalysis populates a ProgressGrid for a group of students, and a particular subject.
func PGAnalysis(subject *Subject, students []Student, nat national.National) ProgressGrid {

	ks2grades := []string{"None", "1", "2", "3C", "3B", "3A", "4C", "4B", "4A", "5C", "5B", "5A", "6"}
	grades := subject.Level.SortedGrades()

	// Initialise grid
	grid := ProgressGrid{KS2: ks2grades, Grades: grades, Cells: map[string](map[string]PGCell){},
		Counts: map[string]int{}, VA: map[string]float64{}}
	for _, ks2 := range ks2grades {
		grid.Cells[ks2] = map[string]PGCell{}
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
		case "Re":
			ks2 = s.KS2.Re
		default:
			ks2 = s.KS2.Av
		}

		if ks2 == "" {
			ks2 = "None"
		}

		cell := grid.Cells[ks2][c.Grd]
		cell.Students = append(cell.Students, PGCellStudent{s.UPN, s.Name()})
		grid.Cells[ks2][c.Grd] = cell
		grid.Counts[c.Grd]++
	}

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

// A PGStudent holds the details of a student and their results, va etc. for populating a table.
type PGStudent struct {
	Student
	Class      string
	KS2        string
	Grade      string
	Effort     int
	VAExists   bool
	VA         float64
	Attendance float64
}

// PGStudentList populates a list of students for the table
func PGStudentList(subject *Subject, students []Student, nat national.National) []PGStudent {

	stdnts := []PGStudent{}
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
		case "Re":
			ks2 = s.KS2.Re
		default:
			ks2 = s.KS2.Av
		}

		va := s.SubjectVA(subject.Subj, nat)

		stdnts = append(stdnts, PGStudent{
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
