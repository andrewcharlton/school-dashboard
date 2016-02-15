package group

import "github.com/andrewcharlton/school-dashboard/analysis/subject"

type ProgressGrid struct {
	Cells   [][]Group
	CellVA  [][]float64
	KS2     []string
	Grades  []string
	VA      []float64
	Counts  []int
	Cohorts []int
}

// ProgressGrid calculates the
func (g Group) ProgressGrid(subject *subject.Subject, natYear string) ProgressGrid {

	cells := [][]Group{}
	cellVA := [][]float64{}
	counts := []int{}
	cohorts := []int{}
	va := []float64{}

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
		va = append(va, 0.0)
		cohorts = append(cohorts, 0)

		for _, g := range grades {
			cells[n] = append(cells[n], Group{})
			tm, exists := subject.TMs[natYear]
			switch exists {
			case true:
				va, err := tm.ValueAdded(ks2, g)
				if err == nil {
					cellVA[n] = append(cellVA[n], va)
				}
			case false:
				cellVA[n] = append(cellVA[n], 0.0)
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
		va[ks2ID] += s.SubjectVA(subject.Subj).Score()
	}

	for n, _ := range ks2Levels {
		if cohorts[n] > 0 {
			va[n] /= float64(cohorts[n])
		}
	}

	return ProgressGrid{Cells: cells, CellVA: cellVA, KS2: ks2Levels, Grades: grades, VA: va, Cohorts: cohorts, Counts: counts}
}
