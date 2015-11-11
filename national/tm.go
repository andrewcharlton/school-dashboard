package national

import (
	"errors"

	"github.com/andrewcharlton/school-dashboard/analysis"
)

// A TMRow holds the probabilities of achieving
// each grade at a level.
type TMRow map[string]float64

// A TransitionMatrix contains the probabilities of achieving
// each grade, at each KS2 starting point.
// It also holds a Level object to enable access to the points
// scores etc. for each grade.
type TransitionMatrix struct {
	Rows  map[string]TMROW
	Level *analysis.Level
}

// Expected calculates the expected Attainment 8 points
func (tm TransitionMatrix) Expected(ks2 string) (float64, error) {

	row, exists := tm.Rows[ks2]
	if !exists {
		return float64(0), errors.New("KS2 score not recognised.")
	}

	total := float64(0)
	for grade, prob := range row.Probs {
		total += tm.Level.Gradeset[grade].Att8 * prob
	}
	return total
}
