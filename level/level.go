// Package level extracts the level struct so it
// can be used in both analysis and national packages.
package level

import "sort"

// A Level brings together the details and grades
// of a particular qualification level.
type Level struct {

	// Name
	Lvl string

	// Does the qualification count as a GCSE
	IsGCSE bool

	// Possible grades achievable at that level
	Gradeset map[string]*Grade
}

// A Grade contains the various points values for a
// grade at a particular level.
type Grade struct {

	// Name of the grade
	Grd string

	// QCA points
	Pts int

	// Attainment 8 points
	Att8 float64

	// Whether the grade counts as a pass at levels 1/2
	L1Pass bool
	L2Pass bool
}

type grdPt struct {
	Grade  string
	Points int
	Att8   float64
}

type grdPts []grdPt

func (g grdPts) Len() int      { return len(g) }
func (g grdPts) Swap(i, j int) { g[i], g[j] = g[j], g[i] }
func (g grdPts) Less(i, j int) bool {
	// Try to sort by points first, otherwise use attainment 8 scores.
	if g[i].Points != g[j].Points {
		return g[i].Points < g[j].Points
	}
	return g[i].Att8 < g[j].Att8
}

// Produces a list of grades, sorted by points score.
func (l Level) SortedGrades() []string {

	grades := grdPts{}
	for _, grade := range l.Gradeset {
		grades = append(grades, grdPt{grade.Grd, grade.Pts, grade.Att8})
	}

	sort.Sort(grades)
	names := []string{}
	for _, g := range grades {
		names = append(names, g.Grade)
	}

	return names
}
