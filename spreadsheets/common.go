// Package spreadsheets produces spreadsheets with the information shown
// on the
package spreadsheets

import (
	"fmt"
	"time"

	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/tealeg/xlsx"
)

func formatBool(b bool) string {

	if b {
		return "Y"
	}
	return "N"
}

func newCell(row *xlsx.Row, contents string, style *xlsx.Style) {

	cell := row.AddCell()
	cell.Value = contents
	cell.SetStyle(style)
}

func newInt(row *xlsx.Row, contents int, style *xlsx.Style) {

	cell := row.AddCell()
	cell.SetInt(contents)
	cell.SetStyle(style)
}

func newFloat(row *xlsx.Row, contents float64, format string, style *xlsx.Style) {

	cell := row.AddCell()
	cell.SetFloatWithFormat(contents, format)
	cell.SetStyle(style)
}

type exportInfo struct {
	Key   string
	Value string
}

// Produces a sheet with all of the export i.nformation held on it.
func exportInfoSheet(sheet *xlsx.Sheet, e database.Env, f database.Filter) {

	d, m, y := time.Now().Date()
	date, _ := e.LookupDate(f.Date)
	rs, _ := e.LookupResultset(f.Resultset)
	nat, _ := e.LookupNatYear(f.NatYear)

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
