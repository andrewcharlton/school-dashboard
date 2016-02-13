package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/analysis/group"
	"github.com/andrewcharlton/school-dashboard/env"
)

// AttendanceGroup summary pages
func AttendanceGroups(e env.Env) http.HandlerFunc {
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

		type YearGroup struct {
			Name   string
			Query  template.URL
			Groups []attendanceSubGroup
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
			[]YearGroup{{"All Years", template.URL(""), attendanceSubGroups(g)}},
		}

		for year := 7; year < 15; year++ {
			y := g.SubGroup(group.Year(year))
			if len(y.Students) == 0 {
				continue
			}
			yeargroup := YearGroup{fmt.Sprintf("Year %v", year),
				template.URL(fmt.Sprintf("&year=%v", year)),
				attendanceSubGroups(y)}
			data.YearGroups = append(data.YearGroups, yeargroup)
		}

		err = e.Templates.ExecuteTemplate(w, "attendancegroups.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}

type attendanceSubGroup struct {
	Name       string
	Query      template.URL
	Attendance group.AttendanceSummary
}

func attendanceSubGroups(g group.Group) []attendanceSubGroup {

	return []attendanceSubGroup{
		{"All", template.URL(""), g.Attendance()},
		{"Male", template.URL("&gender=1"), g.SubGroup(group.Male).Attendance()},
		{"Female", template.URL("&gender=0"), g.SubGroup(group.Female).Attendance()},
		{"Disadvantaged", template.URL("&pp=1"), g.SubGroup(group.PP).Attendance()},
		{"Non-Disadvantaged", template.URL("&pp=0"), g.SubGroup(group.NonPP).Attendance()},
		{"High", template.URL("&ks2band=High"), g.SubGroup(group.High).Attendance()},
		{"Middle", template.URL("&ks2band=Middle"), g.SubGroup(group.Middle).Attendance()},
		{"Low", template.URL("&ks2band=Low"), g.SubGroup(group.Low).Attendance()},
	}
}

// AttendanceExplorer provides a page for exploring the attendance figures
// in more detail, and examine individual students.
func AttendanceExplorer(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if redir := checkRedirect(e, queryOpts{true, true}, w, r); redir {
			return
		}

		Header(e, w, r)
		FilterPage(e, w, r, false)
		defer Footer(e, w, r)

		f := GetFilter(e, r)
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
