package analysis

import (
	"sort"

	"github.com/andrewcharlton/school-dashboard/national"
)

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

func (s slots) Len() int {
	return len(s)
}

func (s slots) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

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
	b := Basket{ks2: s.KS2.APS}

	// Maths baskets
	maths := slots{}
	for _, c := range s.Courses {
		if c.EBacc == "M" {
			maths = append(maths, Slot{c.Subj, c.Grd, c.Att8})
		}
	}
	sort.Sort(sort.Reverse(maths))
	if len(maths) > 0 {
		b.Slots[0] = maths[0]
		b.Slots[1] = maths[0]
		used[maths[0].Subj] = true
	}

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
		b.Slots[2] = eng[0]
		if lang && lit {
			b.Slots[3] = eng[0]
		}
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
	for n, o := range other {
		if n >= 3 {
			break
		}
		b.Slots[n+7] = o
	}

	return b
}

// TotalPoints calculates the total number of points in the basket
func (b Basket) TotalPoints() float64 {

	points := float64(0)
	for _, slot := range b.Slots {
		points += slot.Points
	}
	return points
}

// Entries calculates how many slots in the basket were filled (out of 10)
func (b Basket) Entries() int {

	entries := 0
	for _, slot := range b.Slots {
		if slot.Points > 0 {
			entries++
		}
	}
	return entries
}

// Progress 8 calculates a progress 8 score for the student,
// compared to the national data provided.
func (b Basket) Progress8(nat national.National) Result {

	actual := b.TotalPoints()
	exp, err := nat.Attainment8(b.ks2)
	if err != nil {
		return Result{Error: err}
	}

	return Result{Points: (actual - exp) / 10, Error: nil}
}
