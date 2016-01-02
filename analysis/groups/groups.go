// Package groups provides tools to compare the various subgroups of
//
package groups

import "github.com/andrewcharlton/school-dashboard/analysis"

// GroupDef provides a definition of a group - it's name, and a
// function to determine whether or not a student should be in
// the group or not.
type GroupDef struct {
	Name     string
	Contains func(analysis.Student) bool
}

// Groups contains a list of standard groups for use in comparisons.
var Groups = []GroupDef{
	{"All Students", func(s analysis.Student) bool { return true }},
	{"Male", func(s analysis.Student) bool { return s.Gender == "Male" }},
	{"Female", func(s analysis.Student) bool { return s.Gender == "Female" }},
	{"High", func(s analysis.Student) bool { return s.KS2.Band == "High" }},
	{"Middle", func(s analysis.Student) bool { return s.KS2.Band == "Middle" }},
	{"Low", func(s analysis.Student) bool { return s.KS2.Band == "Low" }},
	{"Disadvantaged", func(s analysis.Student) bool { return s.PP }},
	{"Non-Disadvantaged", func(s analysis.Student) bool { return !s.PP }},
}
