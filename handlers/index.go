package handlers

import (
	"fmt"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/env"
)

// Index produces a landing page
func Index(e env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		header(e, w, r, 0)
		defer footer(e, w, r)

		fmt.Fprintln(w, `<h3>Venerable Bede Data Analysis</h3><br>
					 Please choose an option from the menu bar.<br>
					 <i>Any problems, please let <a href="mailto:andrew.charlton@venerablebede.co.uk">Andrew Charlton</a> know.</i>`)

	}
}
