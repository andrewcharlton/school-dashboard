package stdnt

import "errors"

// EffortAv is a wrapper for average effort
type EffortAv struct {
	Effort float64
	Err    error
}

// Effort calculates the average effort a student.
func (s Student) Effort() EffortAv {

	total, num := 0, 0
	for _, c := range s.Courses {
		total += c.Effort
		num += 1
	}

	if num == 0 {
		return EffortAv{Err: errors.New("No courses present")}
	}
	return EffortAv{Effort: float64(total) / float64(num), Err: nil}
}
