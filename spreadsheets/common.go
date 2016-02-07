// Package spreadsheets produces spreadsheets with the information shown
// on the
package spreadsheets

import (
	"fmt"
	"time"

	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/tealeg/xlsx"
)

var defaultFont = xlsx.Font{11, "Calibri", 2, 0, "FF000000", false, false, false}

type exportInfo struct {
	Key   string
	Value string
}

// Produces a sheet with all of the export i.nformation held on it.
func exportInfoSheet(sheet *xlsx.Sheet, db database.Database, f database.Filter) {

	d, m, y := time.Now().Date()
	date, _ := db.LookupDate(f.Date)
	rs, _ := db.LookupResultset(f.Resultset)
	nat, _ := db.LookupNatYear(f.NatYear)

	info := []exportInfo{{"Export Date: ", fmt.Sprintf("%02d-%02d-%d", d, m, y)},
		{"Effective Date:", date},
		{"Resultset:", rs},
		{"National Data:", nat},
		{"Yeargroup:", f.Year},
	}

	row := sheet.AddRow()
	for _, i := range info {
		row = sheet.AddRow()
		cell := row.AddCell()
		cell = row.AddCell()
		cell.Value = i.Key
		cell = row.AddCell()
		cell.Value = i.Value
		row = sheet.AddRow()
		row.SetHeightCM(0.2)
	}

}
