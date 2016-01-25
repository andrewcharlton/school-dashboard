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
	blankCell(row)
	newCell(row, "Progress Grid: "+subj, newStyle("Title", "None", "None", "Left"))

	row = sheet.AddRow()
	row = sheet.AddRow()
	blankCell(row)
	newCell(row, "KS2", newStyle("Bold", "None", "Bottom", "Left"))
	for _, g := range pg.Grades {
		newCell(row, g, newStyle("Bold", "None", "Bottom", "Center"))
	}
	newCell(row, "VA", newStyle("Bold", "None", "Bottom", "Center"))

	for _, ks2 := range pg.KS2 {
		row := sheet.AddRow()
		blankCell(row)
		newCell(row, ks2, newStyle("Bold", "None", "None", "Left"))
		for _, g := range pg.Grades {
			pgCell := pg.Cells[ks2][g]
			switch {
			case pgCell.VA < -0.33:
				newInt(row, len(pgCell.Students), newStyle("Default", "Red", "None", "Center"))
			case pgCell.VA > 0.67:
				newInt(row, len(pgCell.Students), newStyle("Default", "Green", "None", "Center"))
			default:
				newInt(row, len(pgCell.Students), newStyle("Default", "Yellow", "None", "Center"))
			}

		}
		newFloat(row, pg.VA[ks2], "+0.00;-0.00;0.00", newStyle("Default", "None", "None", "Center"))
	}

	row = sheet.AddRow()
	blankCell(row)
	newCell(row, "Total", newStyle("Bold", "None", "None", "None"))
	for _, g := range pg.Grades {
		newInt(row, pg.Counts[g], newStyle("Bold", "None", "None", "Center"))
	}
	newFloat(row, pg.TotalVA, "+0.00;-0.00;0.00", newStyle("Bold", "None", "None", "Center"))
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
	newCell(row, "Name", newStyle("Bold", "None", "Bottom", "Left"))
	newCell(row, "Class", newStyle("Bold", "None", "Bottom", "Left"))

	headers := []string{"Gender", "PP", "KS2", "SEN", "Grade",
		"Effort", "VA", "Attendance", ""}
	for _, h := range headers {
		newCell(row, h, newStyle("Bold", "None", "Bottom", "Center"))
	}

	// Add historical resultsets to the headers
	for _, rs := range resultsets {
		if !empty[rs.ID] {
			newCell(row, rs.Name, newStyle("Bold", "None", "Bottom", "Vertical"))
		}
	}

	// Add Student data
	for _, s := range studList {
		row := sheet.AddRow()
		newCell(row, s.Name(), newStyle("Default", "None", "None", "Left"))
		newCell(row, s.Class, newStyle("Default", "None", "None", "Left"))
		newCell(row, string(s.Gender[0]), newStyle("Default", "None", "None", "Center"))
		newBool(row, s.PP, newStyle("Default", "None", "None", "Center"))
		newCell(row, s.KS2, newStyle("Default", "None", "None", "Center"))
		newCell(row, s.SEN.Status, newStyle("Default", "None", "None", "Center"))
		newCell(row, s.Grade, newStyle("Default", "None", "None", "Center"))
		newInt(row, s.Effort, newStyle("Default", "None", "None", "Center"))
		newFloat(row, s.VA, "+0.00;-0.00;0.00", newStyle("Default", "None", "None", "Center"))
		newFloat(row, s.Attendance, ".0", newStyle("Default", "None", "None", "Center"))
		newCell(row, "", newStyle("Default", "None", "None", "Center"))
		for _, rs := range resultsets {
			if !empty[rs.ID] {
				grd, _ := historical[rs.ID][s.UPN]
				newCell(row, grd, newStyle("Default", "None", "None", "Center"))
			}
		}

	}

	return nil
}
