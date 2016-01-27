package group

import "errors"

// KS2APS calculates the average points score for
// a group of students.
func (g Group) KS2APS() Result {

	total, num := float64(0), 0
	for _, s := range g.Students {
		if s.KS2.APS != 0 {
			total += s.KS2.APS
			num += 1
		}
	}

	if num == 0 {
		return Result{Pts: 0, Error: errors.New("No KS2 Data available")}
	}
	return Result{Pts: total / float64(num), Error: nil}
}

// KS2Bands returns the number of students in each band
// and what percentage of students is in each.
func (g Group) KS2Bands() map[string]Result {

	count := map[string]int{}
	total := 0
	for _, s := range g.Students {
		count[s.KS2.Band] += 1
		total += 1
	}

	bands := map[string]Result{"High": Result{}, "Middle": Result{},
		"Low": Result{}, "None": Result{}}
	for key, n := range count {
		bands[key] = Result{EntN: n,
			EntP: float64(n) / float64(total)}
	}

	return bands
}
