package handlers

import (
	"fmt"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/env"
)

func Index(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		Header(e, w, r)
		FilterPage(e, w, r, false)

		fmt.Fprintf(w, "Index Page")

		Footer(e, w, r)

	}
}
