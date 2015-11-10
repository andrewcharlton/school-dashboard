package analysis

import "errors"

// Effort calculates the average effort a student.
func (s Student) Effort() Result {

	total, num := 0, 0
	for _, c := range s.Courses {
		total += c.Effort
		num += 1
	}

	if num == 0 {
		return Result{Error: errors.New("No courses present")}
	}

	return Result{EntN: num, Pts: float64(total) / float64(num)}
}
