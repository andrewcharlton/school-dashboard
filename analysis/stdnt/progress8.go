package stdnt

import (
	"sort"

	"github.com/andrewcharlton/school-dashboard/analysis/national"
)

// A P8Score holds the scores for a progress 8 bucket.
type P8Score struct {
	Ent int     // Entries
	Exp float64 // Expected score
	Ach float64 // Achieved score
	Pts float64 // Progress 8 points
	Err error
}

// A Slot holds the details of the subjects and grades
// in the basket.
type Slot struct {
	Subj   string
	Grade  string
	Points float64
}

// slots is used for sorting a selection of subjects
// so that they can be put into the basket.
type slots []Slot

func (s slots) Len() int      { return len(s) }
func (s slots) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s slots) Less(i, j int) bool {

	if s[i].Points == s[j].Points {
		return s[i].Subj > s[j].Subj
	}
	return s[i].Points < s[j].Points
}

// A Basket holds the details of all of the subjects
// that are present in a student's attainment 8 score.
type Basket struct {
	Slots [10]Slot
	ks2   float64
}

// Basket calculates the contents of the students Progress 8
// 'basket', containing Maths, English, 3xEBacc subjects,
// and 3 others.
func (s Student) Basket() Basket {

	used := map[string]bool{}
	b := Basket{Slots: [10]Slot{}, ks2: s.KS2.APS}

	// English baskets - double weighted only if Lang & Lit are
	// present
	lang, lit := false, false
	eng := slots{}
	for _, c := range s.Courses {
		if c.EBacc == "En" {
			eng = append(eng, Slot{c.Subj, c.Grd, c.Att8})
			lang = true
		}
		if c.EBacc == "El" {
			eng = append(eng, Slot{c.Subj, c.Grd, c.Att8})
			lit = true
		}
	}
	sort.Sort(sort.Reverse(eng))
	if len(eng) > 0 {
		used[eng[0].Subj] = true
		b.Slots[0] = eng[0]
		if lang && lit {
			b.Slots[1] = eng[0]
		}
	}

	// Maths baskets
	maths := slots{}
	for _, c := range s.Courses {
		if c.EBacc == "M" {
			maths = append(maths, Slot{c.Subj, c.Grd, c.Att8})
		}
	}
	sort.Sort(sort.Reverse(maths))
	if len(maths) > 0 {
		b.Slots[2] = maths[0]
		b.Slots[3] = maths[0]
		used[maths[0].Subj] = true
	}

	// EBacc Basket
	ebacc := slots{}
	for _, c := range s.Courses {
		if c.EBacc == "H" || c.EBacc == "S" || c.EBacc == "L" {
			ebacc = append(ebacc, Slot{c.Subj, c.Grd, c.Att8})
		}
	}
	sort.Sort(sort.Reverse(ebacc))
	for n, e := range ebacc {
		if n >= 3 {
			break
		}
		used[e.Subj] = true
		b.Slots[n+4] = e
	}

	// Others
	other := slots{}
	for _, c := range s.Courses {
		if !used[c.Subj] {
			other = append(other, Slot{c.Subj, c.Grd, c.Att8})
		}
	}
	sort.Sort(sort.Reverse(other))
	for n, oth := range other {
		if n >= 3 {
			break
		}
		b.Slots[n+7] = oth
	}

	return b
}

// Entries counts the number of slots filled with a non-zero score, up to a maximum
// of 10 slots.
func (b Basket) Entries() int {

	ent, _ := b.points(0, 9)
	return ent
}

// points calculates the entries, number of points in the slots from first to last inclusive.
// first and last should be between 0 and 9.
func (b Basket) points(first, last int) (int, float64) {

	entries := 0
	points := 0.0
	for n := first; n <= last; n++ {
		points += b.Slots[n].Points
		if b.Slots[n].Points > 0 {
			entries++
		}
	}
	return entries, points
}

// Attainment8 calculates the total points achieved in the basket.
func (b Basket) Attainment8() float64 {

	_, pts := b.points(0, 9)
	return pts
}

// Progress 8 calculates a progress 8 score for the student,
// compared to the national data provided.
func (b Basket) Progress8(nat national.Progress8) P8Score {

	entries, actual := b.points(0, 9)
	return P8Score{Ent: entries, Exp: nat.Att8, Ach: actual, Pts: (actual - nat.Att8) / 10.0}
}

// English calculates the progress 8 scores for the English
// element of the Progress 8 Basket.
func (b Basket) English(nat national.Progress8) P8Score {

	entries, actual := b.points(0, 1)
	return P8Score{Ent: entries, Exp: nat.English, Ach: actual, Pts: (actual - nat.English) / 2.0}
}

// Mathematics calculates the progress 8 scores for the Mathematics
// element of the Progress 8 Basket.
func (b Basket) Mathematics(nat national.Progress8) P8Score {

	entries, actual := b.points(2, 3)
	return P8Score{Ent: entries, Exp: nat.Maths, Ach: actual, Pts: (actual - nat.Maths) / 2.0}
}

// EBacc calculates the progress 8 scores for the EBacc
// element of the Progress 8 Basket.
func (b Basket) EBacc(nat national.Progress8) P8Score {

	entries, actual := b.points(4, 6)
	return P8Score{Ent: entries, Exp: nat.EBacc, Ach: actual, Pts: (actual - nat.EBacc) / 3.0}
}

// Other calculates the progress 8 scores for the Other
// element of the Progress 8 Basket.
func (b Basket) Other(nat national.Progress8) P8Score {

	entries, actual := b.points(7, 9)
	return P8Score{Ent: entries, Exp: nat.Other, Ach: actual, Pts: (actual - nat.Other) / 3.0}
}
