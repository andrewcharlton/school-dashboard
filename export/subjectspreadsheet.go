package export

import (
	"database/sql"
	"io"
	"strconv"

	"github.com/andrewcharlton/school-dashboard/analysis/group"
	"github.com/andrewcharlton/school-dashboard/analysis/subject"
	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/andrewcharlton/school-dashboard/env"
	"github.com/tealeg/xlsx"
)

func SubjectSpreadsheet(e env.Env, f database.Filter, subj *subject.Subject, w io.Writer) error {

	file := xlsx.NewFile()

	g, err := e.GroupByFilteredClass(strconv.Itoa(subj.SubjID), "", f)
	if err != nil {
		return err
	}

	sheet, err := file.AddSheet("Progress Grid")
	if err != nil {
		return err
	}
	subjectGrid(sheet, g, subj, f.NatYear)

	sheet, err = file.AddSheet("Students")
	if err != nil {
		return err
	}
	subjectBroadsheet(e, subj, sheet, g)

	sheet, err = file.AddSheet("Export Details")
	if err != nil {
		return err
	}
	exportInfoSheet(sheet, e, f)

	file.Write(w)
	return nil
}

func subjectGrid(sheet *xlsx.Sheet, g group.Group, subj *subject.Subject, natYear string) {

	pg := g.ProgressGrid(subj, natYear)

	row := sheet.AddRow()
	row = sheet.AddRow()
	blankCell(row)
	newCell(row, "Progress Grid: "+subj.Subj, newStyle("Title", "None", "None", "Left"))

	row = sheet.AddRow()
	row = sheet.AddRow()
	blankCell(row)
	newCell(row, "KS2", newStyle("Bold", "None", "Bottom", "Left"))
	for _, grd := range pg.Grades {
		newCell(row, grd, newStyle("Bold", "None", "Bottom", "Center"))
	}
	newCell(row, "VA", newStyle("Bold", "None", "Bottom", "Center"))

	for i, ks2 := range pg.KS2 {
		row := sheet.AddRow()
		blankCell(row)
		newCell(row, ks2, newStyle("Bold", "None", "None", "Left"))
		for j := range pg.Grades {
			pgCell := pg.Cells[i][j]
			switch {
			case pg.CellVA[i][j] < -0.33:
				newInt(row, len(pgCell.Students), newStyle("Default", "Red", "None", "Center"))
			case pg.CellVA[i][j] > 0.67:
				newInt(row, len(pgCell.Students), newStyle("Default", "Green", "None", "Center"))
			default:
				newInt(row, len(pgCell.Students), newStyle("Default", "Yellow", "None", "Center"))
			}

		}
		newFloat(row, pg.RowVA[i], "+0.00;-0.00;0.00", newStyle("Default", "None", "None", "Center"))
	}

	row = sheet.AddRow()
	blankCell(row)
	newCell(row, "Total", newStyle("Bold", "None", "None", "None"))
	for i := range pg.Grades {
		newInt(row, pg.Counts[i], newStyle("Bold", "None", "None", "Center"))
	}
	newFloat(row, g.SubjectVA(subj.Subj).VA, "+0.00;-0.00;0.00", newStyle("Bold", "None", "None", "Center"))
}

func subjectBroadsheet(e env.Env, subj *subject.Subject, sheet *xlsx.Sheet, g group.Group) error {

	// Get a list of all resultsets for historical data
	resultsets := e.Resultsets

	// Create set of maps, keyed by resultsetID, then UPN
	historical := map[string](map[string]string){}
	for _, rs := range resultsets {
		historical[rs.ID] = map[string]string{}
	}

	// Load all historical data from the database
	for _, s := range g.Students {
		grds, err := e.HistoricalResults(s.UPN, subj.SubjID)
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
	newCell(row, "Name", newStyle("Bold", "None", "Bottom", "Left"))
	newCell(row, "Class", newStyle("Bold", "None", "Bottom", "Left"))

	headers := []string{"Gender", "PP", "KS2", "SEN", "Grade",
		"Effort", "VA", "Attendance"}
	for _, h := range headers {
		newCell(row, h, newStyle("Bold", "None", "Bottom", "Center"))
	}

	// Add historical resultsets to the headers
	for _, rs := range resultsets {
		if !empty[rs.ID] {
			newCell(row, rs.Name, newStyle("Bold", "None", "Bottom", "Vertical"))
		}
	}

	for _, h := range []string{"Barriers to Learning", "Intervention"} {
		newCell(row, h, newStyle("Bold", "None", "Bottom", "Center"))
	}

	// Add Student data
	for _, s := range g.Students {
		row := sheet.AddRow()
		newCell(row, s.Name(), newStyle("Default", "None", "None", "Left"))
		newCell(row, s.Class(subj.Subj), newStyle("Default", "None", "None", "Left"))
		newCell(row, s.Gender.String(), newStyle("Default", "None", "None", "Center"))
		newBool(row, s.PP, newStyle("Default", "None", "None", "Center"))
		newCell(row, s.KS2.Score(subj.KS2Prior), newStyle("Default", "None", "None", "Center"))
		newCell(row, s.SEN.Status, newStyle("Default", "None", "None", "Center"))
		newCell(row, s.SubjectGrade(subj.Subj), newStyle("Default", "None", "None", "Center"))
		newCell(row, s.SubjectEffort(subj.Subj), newStyle("Default", "None", "None", "Center"))
		newFloat(row, s.SubjectVA(subj.Subj).Score(), "+0.00;-0.00;0.00", newStyle("Default", "None", "None", "Center"))
		newFloat(row, s.Attendance.Latest(), "0.0%", newStyle("Default", "None", "None", "Center"))
		for _, rs := range resultsets {
			if !empty[rs.ID] {
				grd, _ := historical[rs.ID][s.UPN]
				newCell(row, grd, newStyle("Default", "None", "None", "Center"))
			}
		}

	}

	return nil
}
