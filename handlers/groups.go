package handlers

import (
	"html/template"

	"github.com/andrewcharlton/school-dashboard/analysis/group"
)

type subGroup struct {
	Name  string
	Query template.URL
	Group group.Group
}

// subGroups returns a list of subgroups to be analysed.
func subGroups(g group.Group) []subGroup {

	return []subGroup{
		{"All", template.URL(""), g},
		{"Male", template.URL("&gender=1"), g.SubGroup(group.Male)},
		{"Female", template.URL("&gender=0"), g.SubGroup(group.Female)},
		{"Disadvantaged", template.URL("&pp=1"), g.SubGroup(group.PP)},
		{"Non-Disadvantaged", template.URL("&pp=0"), g.SubGroup(group.NonPP)},
		{"High", template.URL("&ks2band=High"), g.SubGroup(group.High)},
		{"Middle", template.URL("&ks2band=Middle"), g.SubGroup(group.Middle)},
		{"Low", template.URL("&ks2band=Low"), g.SubGroup(group.Low)},
	}
}
