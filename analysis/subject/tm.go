package subject

import "fmt"

// A TransitionMatrix contains the probabilities of achieving
// each grade, at each KS2 starting point.
// Probabilities are indexed by KS2 then Grade
type TransitionMatrix struct {
	Rows map[string](map[string]float64)
	lvl  Level
}

// NewTM creates a new Transition Matrix for a subject at a given level.
func NewTM(lvl Level) TransitionMatrix {

	tm := TransitionMatrix{Rows: map[string](map[string]float64){}, lvl: lvl}
	return tm
}

// Expected calculates the expected Attainment 8 points for a
// given KS2 score.
func (tm TransitionMatrix) Expected(ks2 string) (float64, error) {

	row, exists := tm.Rows[ks2]
	if !exists {
		return 0.0, fmt.Errorf("KS2 score not recognised")
	}

	total := 0.0
	for grade, prob := range row {
		grd, exists := tm.lvl.Gradeset[grade]
		if !exists {
			return 0.0, fmt.Errorf("Grade not recognised in TM: %v", grade)
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
		return 0.0, err
	}

	actual, exists := tm.lvl.Gradeset[grade]
	if !exists {
		return 0.0, fmt.Errorf("Grade not recognised in Level %v - %v", tm.lvl, grade)
	}

	return actual.Att8 - exp, nil
}
