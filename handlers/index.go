package handlers

import (
	"net/http"

	"github.com/andrewcharlton/school-dashboard/database"
)

func Index(e database.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		Header(e, w, r)
		FilterPage(e, w, r, false)

		Footer(e, w, r)

	}
}
