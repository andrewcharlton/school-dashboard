package national

import (
	"errors"

	"github.com/andrewcharlton/school-dashboard/level"
)

// A TMRow holds the probabilities of achieving
// each grade at a level.
type TMRow map[string]float64

// A TransitionMatrix contains the probabilities of achieving
// each grade, at each KS2 starting point.
// It also holds a Level object to enable access to the points
// scores etc. for each grade.
type TransitionMatrix struct {
	Rows  map[string]TMRow
	Level *level.Level
}

// Expected calculates the expected Attainment 8 points
func (tm TransitionMatrix) Expected(ks2 string) (float64, error) {

	row, exists := tm.Rows[ks2]
	if !exists {
		return float64(0), errors.New("KS2 score not recognised.")
	}

	total := float64(0)
	for grade, prob := range row {
		grd, exists := tm.Level.Gradeset[grade]
		if !exists {
			return float64(0), errors.New("Grade not recognised in TM: " + grade)
		}
		total += grd.Att8 * prob
	}
	return total, nil
}

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
