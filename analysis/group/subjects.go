package group

import "fmt"

// SubjectEffort returns the average effort grades recorded in a
// subject.
func (g Group) SubjectEffort(subj string) float64 {

	total, cohort := 0, 0
	for _, s := range g.Students {
		r, exists := s.Results[subj]
		if !exists {
			continue
		}
		total += r.Effort
		cohort++

	}

	if cohort == 0 {
		return 0.0
	}
	return float64(total) / float64(cohort)
}

// SubjectPoints returns the average number of points achieved in
// a subject.
func (g Group) SubjectPoints(subj string) float64 {

	total, cohort := 0, 0
	for _, s := range g.Students {
		r, exists := s.Results[subj]
		if !exists {
			continue
		}
		total += r.Pts
		cohort++
	}

	if cohort == 0 {
		return 0.0
	}
	return float64(total) / float64(cohort)
}

// SubjectVA calculates the overall VA for a group studying a subject.
func (g Group) SubjectVA(subj string) VASummary {

	total := 0.0
	cohort := 0
	for _, s := range g.Students {
		va := s.SubjectVA(subj)
		if va.Err == nil {
			total += va.Score()
			cohort++
		}
	}

	if cohort == 0 {
		return VASummary{0, 0.0, fmt.Errorf("No students with VA scores present.")}
	}
	return VASummary{cohort, total / float64(cohort), nil}
}
