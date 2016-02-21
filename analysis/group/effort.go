package group

// AverageEffort calculates the average effort figure
// for the group in all of their subjects.
func (g Group) AverageEffort() float64 {

	total := 0.0
	n := 0
	for _, s := range g.Students {
		eff := s.AverageEffort()
		if eff.Err == nil {
			total += eff.Effort
			n++
		}
	}

	if n == 0 {
		return 0.0
	}
	return total / float64(n)
}
