package spreadsheets

import (
	"fmt"
	"io"
	"sort"

	"github.com/andrewcharlton/school-dashboard/analysis"
	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/andrewcharlton/school-dashboard/national"
	"github.com/tealeg/xlsx"
)

var ebaccColours = map[string]string{
	"En": "LightBlue",
	"El": "LightBlue",
	"M":  "DarkBlue",
	"S":  "Green",
	"H":  "Green",
	"L":  "Green",
	"":   "Yellow",
}

func Headlines(e database.Env, f database.Filter, w io.Writer) error {

	g, err := e.DB.GroupByFilter(f)
	if err != nil {
		return err
	}

	nat, exists := e.Nationals[f.NatYear]
	if !exists {
		return fmt.Errorf("National year not found: %v", f.NatYear)
	}

	file := xlsx.NewFile()

	sheet, err := file.AddSheet("Broadsheet")
	if err != nil {
		return err
	}
	broadsheet(sheet, g, nat)

	sheet, err = file.AddSheet("Export Details")
	if err != nil {
		return err
	}
	exportInfoSheet(sheet, e, f)

	err = file.Write(w)
	if err != nil {
		return err
	}
	return nil
}

// Produces the broadsheet page
func broadsheet(sheet *xlsx.Sheet, g analysis.Group, nat national.National) error {

	subjects := broadsheetSubjects(g)

	// Create header rows
	err := broadsheetHeaders(sheet, subjects)
	if err != nil {
		return err
	}

	for _, s := range g.Students {
		row := sheet.AddRow()
		broadsheetStudent(row, s, subjects, nat)
	}

	return nil
}

type subject struct {
	Name  string
	EBacc string
}

type subjList []subject

func (s subjList) Len() int      { return len(s) }
func (s subjList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s subjList) Less(i, j int) bool {

	order := map[string]int{
		"En": 0,
		"El": 1,
		"M":  2,
		"S":  3,
		"H":  4,
		"L":  5,
		"":   6,
	}

	switch {
	case order[s[i].EBacc] < order[s[j].EBacc]:
		return true
	case order[s[i].EBacc] > order[s[j].EBacc]:
		return false
	case s[i].Name < s[j].Name:
		return true
	default:
		return false
	}
}

// List all subjects studied by the students.
func broadsheetSubjects(g analysis.Group) subjList {

	// Use a map to create a set of subjects
	subjMap := map[subject]struct{}{}
	for _, s := range g.Students {
		for _, c := range s.Courses {
			header := subject{c.Subj, c.EBacc}
			subjMap[header] = struct{}{}
		}
	}

	// Extract them into a list
	subjects := subjList{}
	for s, _ := range subjMap {
		subjects = append(subjects, s)
	}

	// Sort and return
	sort.Sort(subjects)
	return subjects
}

func broadsheetHeaders(sheet *xlsx.Sheet, subjects subjList) error {

	row := sheet.AddRow()

	// Add Name & usual student filters
	row.SetHeightCM(4.5)
	newCell(row, "UPN", newStyle("Bold", "LightGrey", "Bottom", "Left"))
	newCell(row, "Name", newStyle("Bold", "LightGrey", "Bottom", "Left"))
	for _, f := range studentFilters {
		newCell(row, f, newStyle("Bold", "LightGrey", "Bottom", "Vertical"))
	}
	newCell(row, "", newStyle("Default", "None", "Bottom", "None"))

	for _, s := range subjects {
		newCell(row, s.Name, newStyle("Bold", ebaccColours[s.EBacc], "Bottom", "Vertical"))
	}
	newCell(row, "", newStyle("Default", "None", "Bottom", "None"))

	newCell(row, "Attainment 8", newStyle("Bold", "Orange", "Bottom", "Vertical"))
	newCell(row, "Entries", newStyle("Bold", "Orange", "Bottom", "Vertical"))
	newCell(row, "Progress 8", newStyle("Bold", "Orange", "Bottom", "Vertical"))
	newCell(row, "English", newStyle("Bold", ebaccColours["En"], "Bottom", "Vertical"))
	newCell(row, "Mathematics", newStyle("Bold", ebaccColours["M"], "Bottom", "Vertical"))
	newCell(row, "EBacc", newStyle("Bold", ebaccColours["S"], "Bottom", "Vertical"))
	newCell(row, "Open", newStyle("Bold", ebaccColours[""], "Bottom", "Vertical"))
	newCell(row, "", newStyle("Default", "None", "Bottom", "None"))

	newCell(row, "English & Maths", newStyle("Bold", "Purple", "Bottom", "Vertical"))
	newCell(row, "EBacc: Entered", newStyle("Bold", "Purple", "Bottom", "Vertical"))
	newCell(row, "EBacc: Achieved", newStyle("Bold", "Purple", "Bottom", "Vertical"))
	newCell(row, "", newStyle("Default", "None", "Bottom", "None"))

	newCell(row, "Attendance", newStyle("Bold", "Red", "Bottom", "Vertical"))

	// Hide UPN
	col := sheet.Cols[0]
	col.Hidden = true

	// Make name column wide
	err := sheet.SetColWidth(1, 1, 20.0)
	if err != nil {
		return err
	}

	// Shrink all columns
	err = sheet.SetColWidth(2, sheet.MaxCol, 2.6)
	if err != nil {
		return err
	}
	// Then make Float columns wider
	startCell := len(subjects) + len(studentFilters) + 4
	err = sheet.SetColWidth(startCell, startCell+6, 5.0)
	if err != nil {
		return err
	}
	err = sheet.SetColWidth(startCell+12, startCell+12, 5.0)
	if err != nil {
		return err
	}

	return nil
}

func broadsheetStudent(row *xlsx.Row, s analysis.Student, subjects subjList, nat national.National) {

	studentDetails(row, s)
	blankCell(row)

	for _, subj := range subjects {
		c, exists := s.Courses[subj.Name]
		if exists {
			newCell(row, c.Grd, center)
		} else {
			blankCell(row)
		}
	}
	blankCell(row)

	b := s.Basket()
	newFloat(row, b.Attainment8().Ach, "0.0", center)
	newInt(row, b.Attainment8().EntN, center)

	exp, err := nat.Progress8(s.KS2.APS)
	if err == nil {
		newFloat(row, b.Progress8(exp).Pts, "+0.0;-0.0;0.0", center)
		newFloat(row, b.English(exp).Pts, "+0.0;-0.0;0.0", center)
		newFloat(row, b.Mathematics(exp).Pts, "+0.0;-0.0;0.0", center)
		newFloat(row, b.EBacc(exp).Pts, "+0.0;-0.0;0.0", center)
		newFloat(row, b.Other(exp).Pts, "+0.0;-0.0;0.0", center)
	} else {
		for i := 0; i < 5; i++ {
			blankCell(row)
		}
	}
	blankCell(row)

	newBool(row, s.Basics().AchB, center)
	newBool(row, s.EBacc().Results["All"].EntB, center)
	newBool(row, s.EBacc().Results["All"].AchB, center)
	blankCell(row)

	newFloat(row, s.Attendance.Latest(), "0.0", center)
}
