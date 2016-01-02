package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/andrewcharlton/school-dashboard/analysis/groups"
	"github.com/andrewcharlton/school-dashboard/database"
)

func AttendanceGroups(e database.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		Header(e, w, r)
		FilterPage(e, w, r, true)
		defer Footer(e, w, r)

		f := GetFilter(e, r)

		for year := 7; year <= 11; year++ {
			f.Year = strconv.Itoa(year)
			g, err := e.DB.GroupByFilter(f)
			if err != nil {
				fmt.Fprintf(w, "Error: %v", err)
				return
			}

			summary := groups.Summarise(g.Students, groups.Groups, []groups.ScoreFunc{
				groups.AttendancePct,
				groups.AttendanceUnder(85.0),
				groups.AttendanceUnder(90.0),
			})

			data := struct {
				Title   string
				Summary groups.Summary
			}{
				fmt.Sprintf("Year %v", year),
				summary,
			}

			err = e.Templates.ExecuteTemplate(w, "group.tmpl", data)
			if err != nil {
				fmt.Fprintf(w, "Error: %v", err)
			}
		}
	}
}
