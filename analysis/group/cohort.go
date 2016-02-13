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
		count[s.KS2.Band] += 1
		total += 1
	}

	return count
}
