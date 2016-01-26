// Package national provides data structures
// to hold sets of national data for use
// within the analysis package.
package national

// National holds a set of national data for a year.
type National struct {

	// Attainment 8 point scores for various ks2 scores.
	Prog8 map[string]Progress8

	// Transition Matrices
	TMs map[string]TransitionMatrix
}
