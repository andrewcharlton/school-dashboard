// Awful hacky KS3 work ...
// REMOVE when Life After Levels kicks in.

package student

import "fmt"

// Completely arbitrary sublevel scale
var ks3Levels = map[string]int{
	"1":  1,
	"2":  4,
	"2c": 3,
	"2b": 4,
	"2a": 5,
	"3c": 6,
	"3b": 7,
	"3a": 8,
	"4c": 9,
	"4b": 10,
	"4a": 11,
	"5c": 12,
	"5b": 13,
	"5a": 14,
	"6c": 15,
	"6b": 16,
	"6a": 17,
	"7c": 18,
	"7b": 19,
	"7a": 20,
	"8c": 21,
	"8b": 22,
	"8a": 23,
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

	return VAScore{Expected: float64(exp) / 3.0, Achieved: float64(current-ks2sub) / 3.0, Err: nil}
}
