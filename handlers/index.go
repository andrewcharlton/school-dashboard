package handlers

import (
	"net/http"

	"github.com/andrewcharlton/school-dashboard/env"
)

func Index(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		Header(e, w, r)
		FilterPage(e, w, r, false)

		Footer(e, w, r)

	}
}
