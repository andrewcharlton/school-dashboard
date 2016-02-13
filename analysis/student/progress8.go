package student

import (
	"fmt"
	"sort"
)

// An Attainment8 holds a breakdown of the attainment 8 scores achieved
// in each section of the bucket.
// This is used to store national expectations for students.
type Attainment8 struct {
	English float64
	Maths   float64
	EBacc   float64
	Other   float64
	Overall float64
}

// A Basket holds the contents of a student's Progress 8
// basket.
// Slots are in the order: English (x2), Maths (x2), EBacc (x3), Other (x3)
type Basket struct {
	Slots  [10]Slot
	hasKS2 bool
	nat    Attainment8
}

// A Slot in the basket
type Slot struct {
	Subject string
	Grade   string
	Points  float64
}

// slots is used for sorting a selection of subjects
// so that they can be put into the basket.
type slots []Slot

func (s slots) Len() int           { return len(s) }
func (s slots) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s slots) Less(i, j int) bool { return s[i].Points < s[j].Points }

// Basket returns the contents of a student's basket.
func (s Student) Basket() Basket {

	if s.basket != nil {
		return *s.basket
	}

	basket := Basket{Slots: [10]Slot{}, hasKS2: s.KS2.Exists, nat: s.natAtt8}

	eng := s.engBasket()
	maths := s.mathsBasket()
	ebacc := s.ebaccBasket()

	basket.Slots[0], basket.Slots[1] = eng[0], eng[1]
	basket.Slots[2], basket.Slots[3] = maths[0], maths[1]
	basket.Slots[4], basket.Slots[5], basket.Slots[6] = ebacc[0], ebacc[1], ebacc[2]

	used := map[string]bool{}
	for _, slot := range basket.Slots {
		used[slot.Subject] = true
	}

	other := slots{Slot{}, Slot{}, Slot{}}
	for _, r := range s.Results {
		if used[r.Subj] == false {
			other = append(other, Slot{r.Subj, r.Grd, r.Att8})
		}
	}
	sort.Sort(sort.Reverse(other))
	basket.Slots[7], basket.Slots[8], basket.Slots[9] = other[0], other[1], other[2]

	s.basket = &basket
	return basket
}

func (s Student) engBasket() [2]Slot {

	eng := slots{Slot{}}

	lang, lit := false, false
	for _, r := range s.Results {
		switch r.EBacc {
		case "E":
			lang, lit = true, true
			eng = append(eng, Slot{r.Subj, r.Grd, r.Att8})
		case "En":
			lang = true
			eng = append(eng, Slot{r.Subj, r.Grd, r.Att8})
		case "El":
			lit = true
			eng = append(eng, Slot{r.Subj, r.Grd, r.Att8})
		}
	}

	sort.Sort(sort.Reverse(eng))
	if lang && lit {
		return [2]Slot{eng[0], eng[0]}
	}
	return [2]Slot{eng[0], Slot{}}
}

func (s Student) mathsBasket() [2]Slot {

	maths := slots{Slot{}}
	for _, r := range s.Results {
		if r.EBacc == "M" {
			maths = append(maths, Slot{r.Subj, r.Grd, r.Att8})
		}
	}
	sort.Sort(sort.Reverse(maths))
	return [2]Slot{maths[0], maths[0]}
}

func (s Student) ebaccBasket() [3]Slot {

	ebacc := slots{Slot{}, Slot{}, Slot{}}
	for _, r := range s.Results {
		if r.EBacc == "H" || r.EBacc == "S" || r.EBacc == "L" {
			ebacc = append(ebacc, Slot{r.Subj, r.Grd, r.Att8})
		}
	}
	sort.Sort(sort.Reverse(ebacc))
	return [3]Slot{ebacc[0], ebacc[1], ebacc[2]}
}

// A Progress8Score contains the Attainment and Progress 8 score
// for the student
type Progress8Score struct {
	Entries           int
	Attainment        float64
	AttainmentPerSlot float64
	HasProgress8      bool
	Expected          float64
	Progress8         float64
	Subjects          []string
}

// Basket calculates the Progress 8 score for a particular section
// of the basket.  Start and End are slot numbers, and are inclusive.
func (b Basket) section(start, end int, exp float64) Progress8Score {

	ent := 0
	pts := 0.0
	subjects := []string{}
	for _, s := range b.Slots[start : end+1] {
		switch s.Subject {
		case "":
			subjects = append(subjects, "-")
		default:
			subjects = append(subjects, fmt.Sprintf("%v - %v", s.Subject, s.Grade))
		}
		if s.Points > 0.1 { // 0.1 rather than 0.0 to prevent floating point errors
			ent++
			pts += s.Points
		}
	}

	p8 := Progress8Score{
		Entries:           ent,
		Attainment:        pts,
		AttainmentPerSlot: pts / float64(end+1-start),
		Expected:          exp,
		Subjects:          subjects,
	}
	if exp > 0.1 {
		p8.HasProgress8 = true
		p8.Progress8 = (pts - exp) / float64(end+1-start)
	}
	return p8
}

// Overall produces the progress 8 score for the whole basket
func (b Basket) Overall() Progress8Score { return b.section(0, 9, b.nat.Overall) }

// English produces the progress 8 score the English section of the basket
func (b Basket) English() Progress8Score { return b.section(0, 1, b.nat.English) }

// Maths produces the progress 8 score the Maths section of the basket
func (b Basket) Maths() Progress8Score { return b.section(2, 3, b.nat.Maths) }

// EBacc produces the progress 8 score the EBacc section of the basket
func (b Basket) EBacc() Progress8Score { return b.section(4, 6, b.nat.EBacc) }

// Other produces the progress 8 score the Other section of the basket
func (b Basket) Other() Progress8Score { return b.section(7, 9, b.nat.Other) }
