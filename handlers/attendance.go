package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/analysis/group"
	"github.com/andrewcharlton/school-dashboard/env"
)

// AttendanceGroups produces a page with attendance summaries for the
// various student groups.
func AttendanceGroups(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if redir := checkRedirect(e, w, r, 0); redir {
			return
		}
		header(e, w, r, 0)
		defer footer(e, w, r)

		f := getFilter(e, r)
		g, err := e.GroupByFilter(f)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		type YearGroup struct {
			Name   string
			Query  template.URL
			Groups []subGroup
			Matrix subGroupMatrix
		}

		// Ignore error - will appear as blank string anyway
		week, _ := e.CurrentWeek()

		data := struct {
			Week       string
			Query      template.URL
			YearGroups []YearGroup
		}{
			week,
			template.URL(r.URL.RawQuery),
			[]YearGroup{{"All Years", template.URL(""), subGroups(g), groupMatrix(g)}},
		}

		for year := 7; year < 15; year++ {
			y := g.SubGroup(group.Year(year))
			if len(y.Students) == 0 {
				continue
			}
			yeargroup := YearGroup{fmt.Sprintf("Year %v", year),
				template.URL(fmt.Sprintf("&year=%v", year)),
				subGroups(y),
				groupMatrix(y)}
			data.YearGroups = append(data.YearGroups, yeargroup)

		}

		err = e.Templates.ExecuteTemplate(w, "attendancegroups.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}

// AttendanceExplorer provides a page for exploring the attendance figures
// in more detail, and examine individual students.
func AttendanceExplorer(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if redir := checkRedirect(e, w, r, 2); redir {
			return
		}

		header(e, w, r, 2)
		defer footer(e, w, r)

		f := getFilter(e, r)
		g, err := e.GroupByFilter(f)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		week, _ := e.CurrentWeek()

		data := struct {
			Query template.URL
			Week  string
			Group group.Group
		}{
			template.URL(r.URL.RawQuery),
			week,
			g,
		}

		err = e.Templates.ExecuteTemplate(w, "attendance.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}
