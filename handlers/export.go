package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/andrewcharlton/school-dashboard/env"
	"github.com/andrewcharlton/school-dashboard/export"
)

// ExportSubject
func ExportSubject(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if redir := checkRedirect(e, w, r, 1); redir {
			return
		}
		path := strings.Split(r.URL.Path, "/")
		subjID, err := strconv.Atoi(path[3])
		if err != nil {
			return
		}
		subject := e.Subjects[subjID]

		y, m, d := time.Now().Date()
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		disp := fmt.Sprintf(`attachment; filename="%v %d-%02d-%02d.xlsx"`, subject.Subj, y, m, d)
		w.Header().Set("Content-Disposition", disp)

		f := getFilter(e, r)
		err = export.SubjectSpreadsheet(e, f, subject, w)
		if err != nil {
			fmt.Println(err)
		}
	}
}
