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
	for n, o := range other {
		if n >= 3 {
			break
		}
		b.Slots[n+7] = o
	}

	return b
}

// Attainment 8 calculates the overall attainment 8 score
// for a student.  This is calculated as part of progress 8,
// but this is useful if there is no KS2 data.
func (b Basket) Attainment8() Result {

	actual := 0.0
	entries := 0
	for _, slot := range b.Slots {
		actual += slot.Points
		if slot.Points > 0 {
			entries++
		}
	}
	return Result{Ach: actual, EntN: entries}
}

// Progress 8 calculates a progress 8 score for the student,
// compared to the national data provided.
func (b Basket) Progress8(nat national.Progress8) Result {

	att := b.Attainment8()

	return Result{Ach: att.Ach, Exp: nat.Att8, Pts: (att.Ach - nat.Att8) / 10, EntN: att.EntN}
}

// English calculates the progress 8 scores for the English
// element of the Progress 8 Basket.
func (b Basket) English(nat national.Progress8) Result {

	actual := b.Slots[0].Points + b.Slots[1].Points
	entries := 0
	for _, slot := range b.Slots[2:4] {
		if slot.Points > 0 {
			entries++
		}
	}

	return Result{Ach: actual, Exp: nat.English, Pts: (actual - nat.English) / 2.0, EntN: entries}
}

// Mathematics calculates the progress 8 scores for the Mathematics
// element of the Progress 8 Basket.
func (b Basket) Mathematics(nat national.Progress8) Result {

	actual := b.Slots[2].Points + b.Slots[3].Points
	entries := 0
	for _, slot := range b.Slots[:2] {
		if slot.Points > 0 {
			entries++
		}
	}

	return Result{Ach: actual, Exp: nat.Maths, Pts: (actual - nat.Maths) / 2.0, EntN: entries}
}

// EBacc calculates the progress 8 scores for the EBacc
// element of the Progress 8 Basket.
func (b Basket) EBacc(nat national.Progress8) Result {

	actual := b.Slots[4].Points + b.Slots[5].Points + b.Slots[6].Points
	entries := 0
	for _, slot := range b.Slots[4:7] {
		if slot.Points > 0 {
			entries++
		}
	}

	return Result{Ach: actual, Exp: nat.EBacc, Pts: (actual - nat.EBacc) / 3.0, EntN: entries}
}

// Other calculates the progress 8 scores for the Other
// element of the Progress 8 Basket.
func (b Basket) Other(nat national.Progress8) Result {

	actual := b.Slots[7].Points + b.Slots[8].Points + b.Slots[9].Points
	entries := 0
	for _, slot := range b.Slots[7:10] {
		if slot.Points > 0 {
			entries++
		}
	}

	return Result{Ach: actual, Exp: nat.Other, Pts: (actual - nat.Other) / 3.0, EntN: entries}
}
