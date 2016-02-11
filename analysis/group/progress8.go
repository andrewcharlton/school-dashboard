package group

import "github.com/andrewcharlton/school-dashboard/analysis/student"

// Progress8Summary contains average progress 8 figures for
// the group.  Slots are in order of English, Maths,
// EBacc, Other, Overall.
type Progress8Summary struct {
	Entries           [5]float64
	Attainment        [5]float64
	AttainmentPerSlot [5]float64
	Progress          [5]float64
}

// Progress8 calculates a summary of the progress 8 data
// for a group of students.
func (g Group) Progress8() Progress8Summary {

	cohort := 0
	p8 := Progress8Summary{}
	for _, s := range g.Students {
		b := s.Basket()
		if b.Overall().HasProgress8 {
			cohort++
		}

		// Sum up figures
		for n, score := range map[int]student.Progress8Score{
			0: b.English(),
			1: b.Maths(),
			2: b.EBacc(),
			3: b.Other(),
			4: b.Overall(),
		} {
			p8.Entries[n] += float64(score.Entries)
			p8.Attainment[n] += score.Attainment
			p8.AttainmentPerSlot[n] += score.AttainmentPerSlot
			p8.Progress[n] += score.Progress8
		}
	}

	// Divide to get averages
	for n := 0; n < 5; n++ {
		p8.Entries[n] /= float64(len(g.Students))
		p8.Attainment[n] /= float64(len(g.Students))
		p8.AttainmentPerSlot[n] /= float64(len(g.Students))
		p8.Progress[n] /= float64(cohort)
	}

	return p8
}
