package handlers

import (
	"html/template"

	"github.com/andrewcharlton/school-dashboard/analysis/group"
	"github.com/andrewcharlton/school-dashboard/analysis/student"
)

type groupDef struct {
	Name   string
	Query  string
	Filter (func(s student.Student) bool)
}

var groupDefs = []groupDef{
	{"All", "", func(s student.Student) bool { return true }},
	{"Male", "&gender=1", group.Male},
	{"Female", "&gender=0", group.Female},
	{"Disadvantaged", "&pp=1", group.PP},
	{"Non-Disadvantaged", "&pp=0", group.NonPP},
	{"High", "&ks2band=High", group.High},
	{"Middle", "&ks2band=Middle", group.Middle},
	{"Low", "&ks2band=Low", group.Low},
}

type subGroup struct {
	Name  string
	Query template.URL
	Group group.Group
}

// subGroups returns a list of subgroups to be analysed.
func subGroups(g group.Group) []subGroup {

	groups := []subGroup{}
	for _, def := range groupDefs {
		groups = append(groups, subGroup{def.Name, template.URL(def.Query), g.SubGroup(def.Filter)})
	}
	return groups
}

type subGroupMatrix struct {
	Headers []string
	Groups  [][]subGroup
}

// groupMatrix returns a 2-d list of subgroups.  Headers are used for
// both column and row headers
func groupMatrix(g group.Group) subGroupMatrix {

	headers := []string{}
	var groups [][]subGroup

	for i, def1 := range groupDefs[1:] {
		headers = append(headers, def1.Name)
		row := []subGroup{}
		for j, def2 := range groupDefs[1:] {
			switch {
			case i == j:
				row = append(row, subGroup{def1.Name, template.URL(def1.Query), g.SubGroup(def1.Filter)})
			default:
				row = append(row, subGroup{def1.Name + " & " + def2.Name, template.URL(def1.Query + def2.Query),
					g.SubGroup(def1.Filter, def2.Filter)})
			}
		}
		groups = append(groups, row)
	}

	return subGroupMatrix{headers, groups}
}
