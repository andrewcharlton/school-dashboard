// Awful hacky KS3 work ...
// REMOVE when Life After Levels kicks in.

package student

import "fmt"

// Completely arbitrary sublevel scale
var ks3Levels = map[string]int{
	"1":  1,
	"2":  4,
	"2C": 3,
	"2B": 4,
	"2A": 5,
	"3C": 6,
	"3B": 7,
	"3A": 8,
	"4C": 9,
	"4B": 10,
	"4A": 11,
	"5C": 12,
	"5B": 13,
	"5A": 14,
	"6C": 15,
	"6B": 16,
	"6A": 17,
	"7C": 18,
	"7B": 19,
	"7A": 20,
	"8C": 21,
	"8B": 22,
	"8A": 23,
}

func (s Student) ks3VA(subj string) VAScore {

	r, exists := s.Results[subj]
	if !exists {
		return VAScore{Err: fmt.Errorf("No result for %v in %v found", s.Name(), subj)}
	}

	ks2 := s.KS2.Score(r.KS2Prior)
	if ks2 == "" {
		return VAScore{Err: fmt.Errorf("No %v KS2 score for %v", r.KS2Prior, s.Name())}
	}

	ks2sub, exists := ks3Levels[ks2]
	if !exists {
		return VAScore{Err: fmt.Errorf("KS2 score not recognised: %v", ks2)}
	}

	current, exists := ks3Levels[r.Grd]
	if !exists {
		return VAScore{Err: fmt.Errorf("Current level not recognised: %v", r.Grd)}
	}

	// Hacky hardcoded - 2 sublevels per year, coded for midyear
	var exp int
	switch s.Year {
	case 7:
		exp = 1
	case 8:
		exp = 3
	case 9:
		exp = 5
	}

	return VAScore{Expected: float64(exp), Achieved: float64(current - ks2sub), Err: nil}
}
