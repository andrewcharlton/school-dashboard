package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andrewcharlton/school-dashboard/env"
	"github.com/andrewcharlton/school-dashboard/spreadsheets"
)

func ExportHeadlines(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Produce summary with only year group info
		if redir := checkRedirect(e, w, r, 1); redir {
			return
		}

		y, m, d := time.Now().Date()
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		disp := fmt.Sprintf(`inline; filename="Headlines %d-%02d-%02d.xlsx"`, y, m, d)
		w.Header().Set("Content-Disposition", disp)

		f := getFilter(e, r)

		err := spreadsheets.Headlines(e, f, w)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}
