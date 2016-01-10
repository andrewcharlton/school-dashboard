package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andrewcharlton/school-dashboard/spreadsheets"
)

func ExportSummary(env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Produce summary with only year group info
		if redir := checkRedirect(e, queryOpts{true, false}, w, r); redir {
			return
		}

		y, m, d := time.Now().Date()
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		disp := fmt.Sprintf(`inline; filename="Headlines %d-%02d-%02d.xlsx"`, y, m, d)
		w.Header().Set("Content-Disposition", disp)

		f := GetFilter(e, r)
		g, err := e.DB.GroupByFilter(f)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}
		nat := e.Nationals[f.NatYear]

		err = spreadsheets.Summary(g, f, nat, w)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}
