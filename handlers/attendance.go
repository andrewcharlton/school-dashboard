package handlers

import (
	"fmt"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/analysis/groups"
	"github.com/andrewcharlton/school-dashboard/database"
)

type attData struct {
	Cohort       int
	Possible     int
	Absences     int
	Unauthorised int
	Under85      int
	Under90      int
}

func (a attData) PctAttendance() float64 {

	if a.Possible == 0 {
		return 0.0
	}
	return 100.0 - 100.0*float64(a.Absences)/float64(a.Possible)
}

func Attendance(e database.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if redir := checkRedirect(e, queryOpts{false, false}, w, r); redir {
			return
		}

		Header(e, w, r)
		FilterPage(e, w, r, true)
		defer Footer(e, w, r)

		f := GetFilter(e, r)
		g, err := e.DB.GroupByFilter(f)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		attGroups := map[string]attData{}
		headers := []string{}
		for _, grp := range groups.YearGroups {
			headers = append(headers, grp.Name)
		}

		for _, s := range g.Students {
			for _, grp := range groups.YearGroups {
				if grp.Contains(s) {
					a := attGroups[grp.Name]
					a.Cohort++
					a.Possible += s.Attendance.Possible
					a.Absences += s.Attendance.Absences
					a.Unauthorised += s.Attendance.Unauthorised
					switch {
					case s.Attendance.Latest() < 85.0:
						a.Under85++
						a.Under90++
					case s.Attendance.Latest() < 90.0:
						a.Under90++
					}
					attGroups[grp.Name] = a
				}
			}
		}

		data := struct {
			Headers   []string
			AttGroups map[string]attData
		}{
			headers,
			attGroups,
		}

		err = e.Templates.ExecuteTemplate(w, "attendance.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}
