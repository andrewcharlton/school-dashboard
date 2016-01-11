package spreadsheets

import (
	"database/sql"
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

	sheet, err = file.AddSheet("Students")
	if err != nil {
		return err
	}
	studList := analysis.PGStudentList(subject, g.Students, nat)
	subjPopulateStudents(e, subjID, sheet, studList)

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
		cell.SetInt(pg.Counts[g])
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

func subjPopulateStudents(e database.Env, subjID int, sheet *xlsx.Sheet, studList []analysis.PGStudent) error {

	// Get a list of all resultsets for historical data
	resultsets, err := e.DB.Resultsets()
	if err != nil {
		return err
	}

	// Create set of maps, keyed by resultsetID, then UPN
	historical := map[string](map[string]string){}
	for _, rs := range resultsets {
		historical[rs.ID] = map[string]string{}
	}

	// Load all historical data from the database
	for _, s := range studList {
		grds, err := e.DB.HistoricalResults(s.UPN, strconv.Itoa(subjID))
		switch err {
		case nil:
			for rs, grd := range grds {
				historical[rs][s.UPN] = grd
			}
		case sql.ErrNoRows:
			continue
		default:
			return err
		}
	}

	// Create set of empty resultsets
	empty := map[string]bool{}
	for rs, results := range historical {
		if len(results) == 0 {
			empty[rs] = true
		}
	}

	// Write headers to the sheet
	row := sheet.AddRow()
	row.SetHeightCM(4.5)
	style := xlsx.NewStyle()
	style.Font = boldFont
	style.ApplyFont = true
	style.Border = bottomBorder
	style.ApplyBorder = true
	newCell(row, "Name", style)
	newCell(row, "Class", style)

	style = xlsx.NewStyle()
	style.Font = boldFont
	style.ApplyFont = true
	style.Border = bottomBorder
	style.ApplyBorder = true
	style.Alignment.Horizontal = "center"
	style.ApplyAlignment = true
	headers := []string{"Gender", "PP", "KS2", "SEN", "Grade",
		"Effort", "VA", "Attendance", ""}
	for _, h := range headers {
		newCell(row, h, style)
	}

	// Add historical resultsets to the headers
	style = xlsx.NewStyle()
	style.Font = boldFont
	style.ApplyFont = true
	style.Border = bottomBorder
	style.ApplyBorder = true
	style.Alignment.Horizontal = "center"
	style.Alignment.TextRotation = 90
	style.ApplyAlignment = true
	for _, rs := range resultsets {
		if !empty[rs.ID] {
			newCell(row, rs.Name, style)
		}
	}

	// Add Student data
	style = xlsx.NewStyle()
	style.Font = defaultFont
	style.ApplyFont = true

	center := xlsx.NewStyle()
	center.Font = defaultFont
	center.ApplyFont = true
	center.Alignment.Horizontal = "center"
	center.ApplyAlignment = true
	for _, s := range studList {
		row := sheet.AddRow()
		newCell(row, s.Name(), style)
		newCell(row, s.Class, style)
		newCell(row, string(s.Gender[0]), center)
		newCell(row, formatBool(s.PP), center)
		newCell(row, s.KS2, center)
		newCell(row, s.SEN.Status, center)
		newCell(row, s.Grade, center)
		newInt(row, s.Effort, center)
		newFloat(row, s.VA, "+0.00;-0.00;0.00", center)
		newFloat(row, s.Attendance, ".0", center)
		newCell(row, "", center)
		for _, rs := range resultsets {
			if !empty[rs.ID] {
				grd, _ := historical[rs.ID][s.UPN]
				newCell(row, grd, center)
			}
		}

	}

	return nil
}
