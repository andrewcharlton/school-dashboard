package analysis

import (
	"errors"
	"fmt"

	"github.com/andrewcharlton/school-dashboard/national"
)

// TM retrieves the transition matrix for a particular subject.
func (s Student) TM(subject string, nat national.National) (national.TransitionMatrix, error) {

	c, exists := s.Courses[subject]
	if !exists {
		return national.TransitionMatrix{},
			fmt.Errorf("Course not recognised. UPN: %s, Subject: %v", s.UPN, subject)
	}

	tm, exists := nat.TMs[c.TM]
	if !exists {
		return national.TransitionMatrix{}, fmt.Errorf("TM not found: %v", c.TM)
	}

	return tm, nil
}

// SubjectVA calculates the value added score for a student in a particular
// subject.
func (s Student) SubjectVA(subject string, nat national.National) Result {

	c, exists := s.Courses[subject]
	if !exists {
		return Result{Error: errors.New("No results for subject found")}
	}

	if c.Grd == "" {
		return Result{Error: errors.New("No grade found")}
	}

	tm, exists := nat.TMs[c.TM]
	if !exists {
		return Result{Error: errors.New("No recognised TM")}
	}

	var exp float64
	var err error
	switch c.KS2Prior {
	case "En":
		exp, err = tm.Expected(s.KS2.En)
	case "Ma":
		exp, err = tm.Expected(s.KS2.Ma)
	default:
		exp, err = tm.Expected(s.KS2.Av)
	}
	if err != nil {
		return Result{Error: err}
	}

	return Result{Pts: c.Att8 - exp, Exp: exp}
}
