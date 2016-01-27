package national

import "errors"

// ValueAdded calculates the value added score for an individual
// student, based on their KS2 and Grade achieved.
func (tm TransitionMatrix) ValueAdded(ks2, grade string) (float64, error) {

	exp, err := tm.Expected(ks2)
	if err != nil {
		return float64(0), err
	}

	actual, exists := tm.Level.Gradeset[grade]
	if !exists {
		return float64(0), errors.New("Grade not recognised in TM: " + grade)
	}

	return actual.Att8 - exp, nil
}
