package national

import (
	"errors"
	"fmt"
)

// Progress8 holds the expected values in each area
// of the progress 8 basket.
type Progress8 struct {
	English float64
	Maths   float64
	EBacc   float64
	Other   float64
	Att8    float64
}

var (
	ErrNoKS2       = errors.New("No KS2 Data found.")
	ErrKS2NotFound = errors.New("KS2 score not found in national lookup.")
)

// Progress8 returns the expected Attainment 8 points score
// for a student with the given KS2 APS score.
func (n National) Progress8(ks2aps float64) (Progress8, error) {

	if ks2aps == 0 {
		return Progress8{}, ErrNoKS2
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
		return Progress8{}, ErrKS2NotFound
	}

	return prog8, nil
}
