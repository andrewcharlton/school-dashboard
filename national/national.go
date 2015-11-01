// Package national provides data structures
// to hold sets of national data for use
// within the analysis package.
package national

import "fmt"

// National holds a set of national data for a year.
type National struct {

	// Attainment 8 point scores for various ks2 scores.
	Att8 map[string]float64
}

// Attainment8 returns the expected Attainment 8 points score
// for a student with the given KS2 APS score.
func (n National) Attainment8(ks2aps float64) (float64, error) {

	if ks2aps == 0 {
		return 0, fmt.Errorf("No KS2 data")
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

	att8, exists := n.Att8[dec]
	if !exists {
		return 0, fmt.Errorf("Decimal level not found in table: %v", dec)
	}

	return att8, nil
}
