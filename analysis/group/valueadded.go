package group

// AverageVA calculates the average value added for the group
func (g Group) AverageVA() float64 {

	total := 0.0
	n := 0

	for _, s := range g.Students {
		va := s.AverageVA()
		if va.Err == nil {
			total += va.Score()
			n++
		}
	}

	if n == 0 {
		return 0.0
	}

	return total / float64(n)
}
