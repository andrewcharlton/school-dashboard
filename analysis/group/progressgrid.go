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
		if va.Err != nil {
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
			// Ignore error handling because value defaults to 0.0 anyway
			tm, _ := subject.TMs[natYear]
			va, _ := tm.ValueAdded(ks2, g)
			cellVA[n] = append(cellVA[n], va)
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
