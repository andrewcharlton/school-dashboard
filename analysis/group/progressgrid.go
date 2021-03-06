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

	ks2Levels := []string{"None", "1", "2", "3c", "3b", "3a", "4c", "4b", "4a", "5c", "5b", "5a", "6"}
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

	for n := range ks2Levels {
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
	"2c": 3,
	"2b": 4,
	"2a": 5,
	"3c": 6,
	"3b": 7,
	"3a": 8,
	"4c": 9,
	"4b": 10,
	"4a": 11,
	"5c": 12,
	"5b": 13,
	"5a": 14,
	"6c": 15,
	"6b": 16,
	"6a": 17,
	"7c": 18,
	"7b": 19,
	"7a": 20,
	"8c": 21,
	"8b": 22,
	"8a": 23,
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

	return float64(currSub-ks2Sub-2*year+13) / 3.0, nil
}
