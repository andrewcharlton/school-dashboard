package handlers

import (
	"fmt"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/analysis/group"
	"github.com/andrewcharlton/school-dashboard/env"
)

// Attendance summary pages
func Attendance(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if redir := checkRedirect(e, queryOpts{false, false}, w, r); redir {
			return
		}

		Header(e, w, r)
		FilterPage(e, w, r, true)
		defer Footer(e, w, r)

		f := GetFilter(e, r)
		g, err := e.GroupByFilter(f)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		type SubGroup struct {
			Name       string
			Attendance group.AttendanceSummary
		}

		type YearGroup struct {
			Year   int
			Groups []SubGroup
		}

		data := struct {
			All        group.AttendanceSummary
			YearGroups []YearGroup
		}{
			g.Attendance(),
			[]YearGroup{},
		}

		for year := 7; year < 15; year++ {
			y := g.SubGroup(group.Year(year))
			if len(y.Students) == 0 {
				break
			}
			yeargroup := YearGroup{year, []SubGroup{
				{"All", y.Attendance()},
				{"Male", y.SubGroup(group.Male).Attendance()},
				{"Female", y.SubGroup(group.Female).Attendance()},
				{"Disadvantaged", y.SubGroup(group.PP).Attendance()},
				{"Non-Disadvantaged", y.SubGroup(group.NonPP).Attendance()},
				{"High", y.SubGroup(group.High).Attendance()},
				{"Middle", y.SubGroup(group.Middle).Attendance()},
				{"Low", y.SubGroup(group.Low).Attendance()},
			},
			}
			data.YearGroups = append(data.YearGroups, yeargroup)
		}

		err = e.Templates.ExecuteTemplate(w, "attendance.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}
