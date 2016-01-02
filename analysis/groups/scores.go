package groups

// A Score contains the results of a student's achievement in some measure
// and whether that measure was applicable to them.
type Score struct {
	Score float64
	Error bool
}

// A groupScore collates the scores for a particular measure together.
type groupScore []Score

// Mean calculates the arithmetic mean of any values in the
// groupscore, ignoring any errors.
func (g groupScore) Mean() Score {

	total, n := 0.0, 0.0
	for _, s := range g {
		if !s.Error {
			total += s.Score
			n++
		}
	}

	if n == 0.0 {
		return Score{0.0, true}
	}
	return Score{total / n, false}
}
