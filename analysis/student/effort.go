package student

import "fmt"

// EffortAv wraps the average effort figures for use in templates.
type EffortAv struct {
	Effort float64
	Err    error
}

// Effort calculates the average effort a student has achieved
// in their subjects.
func (s Student) AverageEffort() EffortAv {

	total, num := 0, 0
	for _, r := range s.Results {
		total += r.Effort
		num++
	}

	if num == 0 {
		return EffortAv{0.0, fmt.Errorf("No courses present")}
	}
	return EffortAv{Effort: float64(total) / float64(num), Err: nil}
}

// SubjectEffort provides a wrapper to look up a student's effort
// in a subject for use in templates.  Returns "" if no effort found.
func (s Student) SubjectEffort(subj string) string {

	r, exists := s.Results[subj]
	if !exists {
		return ""
	}
	return fmt.Sprintf("%v", r.Effort)
}
