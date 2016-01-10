// Package exports contains functions to produce spreadsheets
// exporting the data to
package spreadsheets

import (
	"fmt"
	"io"
	"time"

	"github.com/andrewcharlton/school-dashboard/analysis"
	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/andrewcharlton/school-dashboard/national"
	"github.com/tealeg/xlsx"
)

func Summary(g analysis.Group, f database.Filter, nat national.National, w io.Writer) error {

	file := xlsx.NewFile()

	sheet, err := file.AddSheet("Summary")
	if err != nil {
		return err
	}
	err = addExportInfo(sheet, f)
	if err != nil {
		return err
	}

	err = file.Write(w)
	if err != nil {
		return err
	}
	return nil
}

func addExportInfo(sheet *xlsx.Sheet, f database.Filter) error {

	row := sheet.AddRow()
	row = sheet.AddRow()
	cell := row.AddCell()
	cell = row.AddCell()
	cell.Value = "Export Date:"
	cell = row.AddCell()
	y, m, d := time.Now().Date()
	cell.Value = fmt.Sprintf("%v-%v-%v", d, m, y)

}
