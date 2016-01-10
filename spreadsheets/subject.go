package spreadsheets

import (
	"io"
	"strconv"

	"github.com/andrewcharlton/school-dashboard/analysis"
	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/tealeg/xlsx"
)

func Subject(e database.Env, f database.Filter, subjID int, w io.Writer) error {

	file := xlsx.NewFile()

	subject := e.DB.Subjects()[subjID]
	nat := e.Nationals[f.NatYear]
	g, err := e.DB.GroupByFilteredClass(strconv.Itoa(subjID), "", f)
	if err != nil {
		return err
	}

	sheet, err := file.AddSheet("Progress Grid")
	if err != nil {
		return err
	}
	pg := analysis.PGAnalysis(subject, g.Students, nat)
	subjPopulateGrid(sheet, pg, subject.Subj)

	file.Write(w)
	return nil
}

func subjPopulateGrid(sheet *xlsx.Sheet, pg analysis.ProgressGrid, subj string) {

	row := sheet.AddRow()
	row = sheet.AddRow()
	cell := row.AddCell()
	cell = row.AddCell()
	cell.Value = "Progress Grid: " + subj
	style := xlsx.NewStyle()
	style.Font = titleFont
	style.ApplyFont = true
	cell.SetStyle(style)
	row = sheet.AddRow()

	row = sheet.AddRow()
	cell = row.AddCell()
	cell = row.AddCell()
	cell.Value = "KS2"
	style = xlsx.NewStyle()
	style.Font = boldFont
	style.Border = bottomBorder
	style.ApplyBorder = true
	style.ApplyFont = true
	cell.SetStyle(style)

	style.Alignment.Horizontal = "center"
	style.ApplyAlignment = true
	grades := append(pg.Grades, "VA")
	for _, g := range grades {
		cell := row.AddCell()
		cell.Value = g
		cell.SetStyle(style)
	}

	for _, ks2 := range pg.KS2 {
		row := sheet.AddRow()
		cell := row.AddCell()
		cell = row.AddCell()
		cell.Value = ks2
		style := xlsx.NewStyle()
		style.Font = boldFont
		style.ApplyFont = true
		cell.SetStyle(style)

		for _, g := range pg.Grades {
			pgCell := pg.Cells[ks2][g]
			cell := row.AddCell()
			cell.SetInt(len(pgCell.Students))
			style := xlsx.NewStyle()
			style.Font = defaultFont
			style.ApplyFont = true
			style.Border = allBorders
			style.ApplyBorder = true
			style.Alignment.Horizontal = "center"
			style.ApplyAlignment = true
			switch {
			case pgCell.VA < -0.33:
				style.Fill.FgColor = "FFF7774F"
			case pgCell.VA > 0.67:
				style.Fill.FgColor = "FF2FED4F"
			default:
				style.Fill.FgColor = "FFF2EE54"
			}
			style.Fill.PatternType = "solid"
			style.ApplyFill = true
			cell.SetStyle(style)
		}

		cell = row.AddCell()
		cell.SetFloat(pg.VA[ks2])
		cell.NumFmt = "+0.00;-0.00;0.00"
		style = xlsx.NewStyle()
		style.Alignment.Horizontal = "center"
		style.ApplyAlignment = true
		cell.SetStyle(style)
	}

	row = sheet.AddRow()
	cell = row.AddCell()
	cell = row.AddCell()
	cell.Value = "Total"
	style = xlsx.NewStyle()
	style.Font = boldFont
	style.ApplyFont = true
	cell.SetStyle(style)

	for _, g := range pg.Grades {
		cell := row.AddCell()
		cell.SetFloat(pg.VA[g])
		cell.NumFmt = "+0.00;-0.00;0.00"
		style := xlsx.NewStyle()
		style.Font = boldFont
		style.ApplyFont = true
		style.Alignment.Horizontal = "center"
		style.ApplyAlignment = true
		cell.SetStyle(style)
	}

	cell = row.AddCell()
	cell.SetFloat(pg.TotalVA)
	cell.NumFmt = "+0.00;-0.00;0.00"
	style = xlsx.NewStyle()
	style.Font = boldFont
	style.ApplyFont = true
	style.Alignment.Horizontal = "center"
	style.ApplyAlignment = true
	cell.SetStyle(style)

}
