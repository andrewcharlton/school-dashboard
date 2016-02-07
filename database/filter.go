package database

import (
	"fmt"
	"strings"
)

// A Filter holds all the fields available to select a group
// of students on.
type Filter struct {
	// Effective date to use for membership of the group.
	// Holds the id number, but held as string for easier
	// parsing from query strings.
	Date string
	// Which set of assessment results to use.
	// Holds id number as above.
	Resultset string
	// Which year to use for national data comparisons.
	NatYear string
	// Which yeargroup to include in the group
	Year string
	// Pupil premium: "", "1", or "0" for Any/True/False
	PP string
	// EAL students: "", "1", or "0" for Any/True/False
	EAL string
	// Gender: "", "1", "0" for Any/Male/Female
	Gender string
	// SEN types to include - Empty for any
	SEN []string
	// Which Ethnic groups to include, empty for any
	Ethnicities []string
	// Which KS2 bands to include, empty for any
	KS2Bands []string
}

// sql generates a query string from the filter, used to select
// student upns.
func (f Filter) sql(table string) string {

	query := fmt.Sprintf(`SELECT upn FROM %v
						  WHERE date_id=%v`,
		table, f.Date)

	if f.Year != "" {
		query += fmt.Sprintf(" AND year=%v", f.Year)
	}
	if f.PP != "" {
		query += fmt.Sprintf(" AND pp=%v", f.PP)
	}
	if f.EAL != "" {
		query += fmt.Sprintf(" AND eal=%v", f.EAL)
	}
	if f.Gender != "" {
		query += fmt.Sprintf(" AND gender=%v", f.Gender)
	}
	if len(f.SEN) > 0 {
		query += fmt.Sprintf(` AND sen_status IN ("` + strings.Join(f.SEN, `", "`) + `")`)
	}
	if len(f.Ethnicities) > 0 {
		query += fmt.Sprintf(` AND ethnicity IN ("` + strings.Join(f.Ethnicities, `", "`) + `")`)
	}
	if len(f.KS2Bands) > 0 {
		query += fmt.Sprintf(` AND ks2_band IN ("` + strings.Join(f.KS2Bands, `", "`) + `")`)
	}
	return query
}
