package group

// Cohort returns the size of the cohort
func (g Group) Cohort() int {
	return len(g.Students)
}

// KS2APS calculates the average points score for
// a group of students.
func (g Group) KS2APS() float64 {

	total, num := 0.0, 0
	for _, s := range g.Students {
		if s.KS2.APS != 0.0 {
			total += s.KS2.APS
			num++
		}
	}

	if num == 0 {
		return 0.0
	}
	return total / float64(num)
}

// KS2Bands returns the number of students in each band
func (g Group) KS2Bands() map[string]int {

	count := map[string]int{}
	total := 0
	for _, s := range g.Students {
		count[s.KS2.Band]++
		total++
	}

	return count
}

// PP returns the percentage of students who are in
// receipt of pupil premium
func (g Group) PP() float64 {

	if len(g.Students) == 0 {
		return 0.0
	}

	pp := 0
	for _, s := range g.Students {
		if s.PP {
			pp++
		}
	}

	return float64(pp) / float64(len(g.Students))
}
