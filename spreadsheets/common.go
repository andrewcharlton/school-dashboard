// Package spreadsheets produces spreadsheets with the information shown
// on the
package spreadsheets

import (
	"fmt"
	"time"

	"github.com/andrewcharlton/school-dashboard/analysis"
	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/tealeg/xlsx"
)

func blankCell(row *xlsx.Row) {
	row.AddCell()
}

func newBool(row *xlsx.Row, contents bool, style *xlsx.Style) {

	cell := row.AddCell()
	switch contents {
	case true:
		cell.Value = "Y"
	case false:
		cell.Value = "N"
	}
	cell.SetStyle(style)
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
		blankCell(row)
		newCell(row, i.Key, newStyle("Bold", "None", "None", "Left"))
		newCell(row, i.Value, newStyle("Bold", "None", "None", "Left"))
		row.SetHeightCM(0.2)
	}

	sheet.SetColWidth(1, 1, 14.0)

}

// A list of student filters
var studentFilters = []string{
	"Gender",
	"PP",
	"KS2",
	"SEN",
}

// Writes the first few columns to match the student filters
func studentDetails(row *xlsx.Row, s analysis.Student) {

	newCell(row, s.UPN, left)
	newCell(row, s.Name(), left)
	newCell(row, string(s.Gender[0]), center)
	newBool(row, s.PP, center)
	newCell(row, s.KS2.Av, center)
	newCell(row, s.SEN.Status, center)
}
