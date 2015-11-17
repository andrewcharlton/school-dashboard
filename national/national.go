// Package national provides data structures
// to hold sets of national data for use
// within the analysis package.
package national

import "fmt"

// National holds a set of national data for a year.
type National struct {

	// Attainment 8 point scores for various ks2 scores.
	Prog8 map[string]Progress8

	// Transition Matrices
	TMs map[string]TransitionMatrix
}

// Attainment8 returns the expected Attainment 8 points score
// for a student with the given KS2 APS score.
func (n National) Progress8(ks2aps float64) (Progress8, error) {

	if ks2aps == 0 {
		return Progress8{}, fmt.Errorf("No KS2 data")
	}

	var dec string
	switch {
	case ks2aps < 1.55*6:
		dec = "1.5"
	case ks2aps < 2.05*6:
		dec = "2.0"
	case ks2aps < 2.55*6:
		dec = "2.5"
	case ks2aps < 2.85*6:
		dec = "2.8"
	case ks2aps >= 5.75*6:
		dec = "5.8"
	default:
		dec = fmt.Sprintf("%1.1f", ks2aps/6)
	}

	prog8, exists := n.Prog8[dec]
	if !exists {
		return Progress8{}, fmt.Errorf("Decimal level not found in table: %v", dec)
	}

	return prog8, nil
}
