// Package groups provides tools to compare the various subgroups of
//
package groups

import (
	"github.com/andrewcharlton/school-dashboard/analysis/stdnt"
)

// GroupDef provides a definition of a group - it's name, and a
// function to determine whether or not a student should be in
// the group or not.
type GroupDef struct {
	Name     string
	Contains func(analysis.Student) bool
}

// Groups contains a list of standard groups for use in comparisons.
var Groups = []GroupDef{
	{"All Students", func(s stdnt.Student) bool { return true }},
	{"Male", func(s stdnt.Student) bool { return s.Gender == "Male" }},
	{"Female", func(s stdnt.Student) bool { return s.Gender == "Female" }},
	{"High", func(s stdnt.Student) bool { return s.KS2.Band == "High" }},
	{"Middle", func(s stdnt.Student) bool { return s.KS2.Band == "Middle" }},
	{"Low", func(s stdnt.Student) bool { return s.KS2.Band == "Low" }},
	{"Disadvantaged", func(s stdnt.Student) bool { return s.PP }},
	{"Non-Disadvantaged", func(s stdnt.Student) bool { return !s.PP }},
}

var YearGroups = []GroupDef{
	{"All Students", func(s stdnt.Student) bool { return true }},
	{"Year 7", func(s stdnt.Student) bool { return s.Year == 7 }},
	{"Year 8", func(s stdnt.Student) bool { return s.Year == 8 }},
	{"Year 9", func(s stdnt.Student) bool { return s.Year == 9 }},
	{"Year 10", func(s stdnt.Student) bool { return s.Year == 10 }},
	{"Year 11", func(s stdnt.Student) bool { return s.Year == 11 }},
}
