package group

import (
	"fmt"

	"github.com/andrewcharlton/school-dashboard/analysis/subject"
)

// A VASummary holds the details of
type VASummary struct {
	Cohort int
	VA     float64
	Err    error
}

// SubjectVA calculates the overall VA for a group studying a subject.
func (g Group) SubjectVA(subj string) VASummary {

	total := 0.0
	cohort := 0
	for _, s := range g.Students {
		va := s.SubjectVA(subj)
		if va.Err == nil {
			total += va.Score()
			cohort++
		}
	}

	if cohort == 0 {
		return VASummary{0, 0.0, fmt.Errorf("No students with VA scores present.")}
	}
	return VASummary{cohort, total / float64(cohort), nil}
}

// A ProgressGrid holds all of the details for filling out a progress grid
// for the subject.
type ProgressGrid struct {
	Cells   [][]Group   // A list of students at each combination of KS2 & Grade
	CellVA  [][]float64 // The VA score for students each cell
	KS2     []string    // A sorted list of valid KS2 scores
	Grades  []string    // A sorted list of valid grades
	RowVA   []float64   // A list of the total VA score for each row
	Counts  []int       // A count of the number of students achieving each grade
	Cohorts []int       // A list of the number of students present at each KS2 level
}

// ProgressGrid calculates the
func (g Group) ProgressGrid(subject *subject.Subject, natYear string) ProgressGrid {

	cells := [][]Group{}
	cellVA := [][]float64{}
	counts := []int{}
	cohorts := []int{}
	rowVA := []float64{}

	grades := subject.Level.SortedGrades()
	gradeMap := map[string]int{}
	for n, g := range grades {
		gradeMap[g] = n
		counts = append(counts, 0)
	}

	// KS3! Remove when no longer needed.
	var year int
	if len(g.Students) > 0 {
		year = g.Students[0].Year
	}

	ks2Levels := []string{"None", "1", "2", "3C", "3B", "3A", "4C", "4B", "4A", "5C", "5B", "5A", "6"}
	ks2Map := map[string]int{}
	for n, ks2 := range ks2Levels {
		ks2Map[ks2] = n
		cells = append(cells, []Group{})
		cellVA = append(cellVA, []float64{})
		rowVA = append(rowVA, 0.0)
		cohorts = append(cohorts, 0)

		for _, g := range grades {
			cells[n] = append(cells[n], Group{})

			// KS3! Remove when no longer needed
			switch year {
			case 7, 8, 9:
				va, _ := ks3VA(ks2, g, year)
				cellVA[n] = append(cellVA[n], va)
			default:
				// Ignore error handling because value defaults to 0.0 anyway
				tm, _ := subject.TMs[natYear]
				va, _ := tm.ValueAdded(ks2, g)
				cellVA[n] = append(cellVA[n], va)
			}
		}
	}

	for _, s := range g.Students {
		r, exists := s.Results[subject.Subj]
		if !exists {
			continue
		}

		ks2 := s.KS2.Score(subject.KS2Prior)
		if ks2 == "" {
			ks2 = "None"
		}

		ks2ID := ks2Map[ks2]
		grdID := gradeMap[r.Grd]
		cells[ks2ID][grdID].Students = append(cells[ks2ID][grdID].Students, s)
		counts[grdID]++
		cohorts[ks2ID]++
		rowVA[ks2ID] += s.SubjectVA(subject.Subj).Score()
	}

	for n, _ := range ks2Levels {
		if cohorts[n] > 0 {
			rowVA[n] /= float64(cohorts[n])
		}
	}

	return ProgressGrid{
		Cells:   cells,
		CellVA:  cellVA,
		KS2:     ks2Levels,
		Grades:  grades,
		RowVA:   rowVA,
		Counts:  counts,
		Cohorts: cohorts,
	}
}

// Completely arbitrary sublevel scale
var ks3Levels = map[string]int{
	"1":  1,
	"2":  4,
	"3":  7,
	"4":  10,
	"5":  13,
	"6":  16,
	"2C": 3,
	"2B": 4,
	"2A": 5,
	"3C": 6,
	"3B": 7,
	"3A": 8,
	"4C": 9,
	"4B": 10,
	"4A": 11,
	"5C": 12,
	"5B": 13,
	"5A": 14,
	"6C": 15,
	"6B": 16,
	"6A": 17,
	"7C": 18,
	"7B": 19,
	"7A": 20,
	"8C": 21,
	"8B": 22,
	"8A": 23,
}

// Hacky hard-coded level scale
func ks3VA(ks2, current string, year int) (float64, error) {

	ks2Sub, exists := ks3Levels[ks2]
	if !exists {
		return 0.0, fmt.Errorf("KS2 not recognised: %v", ks2)
	}

	currSub, exists := ks3Levels[current]
	if !exists {
		return 0.0, fmt.Errorf("Current level not recognised: %v", current)
	}

	return float64(currSub - ks2Sub - 2*year + 13), nil
}
