package handlers

import (
	"fmt"
	"net/http"

	"github.com/andrewcharlton/school-dashboard/database"
)

func Index(e database.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		Header(e, w, r)
		FilterPage(e, w, r, false)
		defer Footer(e, w, r)

		fmt.Fprintln(w, `<h3>Venerable Bede Data Analysis</h3><br>
					 Please choose an option from the menu bar.<br>
					 <i>Any problems, please let <a href="mailto:andrew.charlton@venerablebede.co.uk">Andrew Charlton</a> know.</i>`)

	}
}
