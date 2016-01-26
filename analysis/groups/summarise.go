package groups

import "github.com/andrewcharlton/school-dashboard/analysis/stdnt"

// A GroupSummary provides a summary of the group and the average for each
// measure.
type GroupSummary struct {
	Name   string
	Cohort int
	Scores []Score
}

// A Summary contains the results of a
type Summary struct {
	Headers []string
	Groups  []GroupSummary
}

// Summarise finds the average score for each subgroup of students (as defined by the
// group definitions), for each measure defined by the ScoreFuncs.
func Summarise(stdnts []stdnt.Student, grps []GroupDef, funcs []ScoreFunc) Summary {

	summary := Summary{Headers: []string{}, Groups: []GroupSummary{}}
	for _, f := range funcs {
		summary.Headers = append(summary.Headers, f.Name)
	}

	// data - indexed first by group name, then measure
	data := map[string](map[string]groupScore){}
	for _, g := range grps {
		data[g.Name] = map[string]groupScore{}
	}
	cohorts := map[string]int{}

	for _, s := range stdnts {
		// Calculate which groups the student is in.
		groups := []string{}
		for _, g := range grps {
			if g.Contains(s) {
				groups = append(groups, g.Name)
				cohorts[g.Name]++
			}
		}

		// Calculate a student's score for each measure and then add it
		// to each list of the group
		for _, f := range funcs {
			score := f.Func(s)
			for _, g := range groups {
				data[g][f.Name] = append(data[g][f.Name], score)
			}
		}
	}

	// Summarise data
	for _, g := range grps {
		grpSum := GroupSummary{Name: g.Name, Cohort: cohorts[g.Name], Scores: []Score{}}
		for _, f := range funcs {
			grpSum.Scores = append(grpSum.Scores, data[g.Name][f.Name].Mean())
		}
		summary.Groups = append(summary.Groups, grpSum)
	}

	return summary
}
