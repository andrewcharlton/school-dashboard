package handlers

import (
	"fmt"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/database"
)

func Headlines(e database.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		Header(e, w, r)
		FilterPage(e, w, r, false)

		f := GetFilter(e, r)
		g, err := e.DB.GroupByFilter(f)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}

		fmt.Fprintf(w, "Basics: %3.1f", float64(100)*g.Basics().AchP)

		Footer(e, w, r)
	}
}
