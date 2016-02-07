package spreadsheets

import (
	"io"
	"sort"

	"github.com/andrewcharlton/school-dashboard/analysis/group"
	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/tealeg/xlsx"
)

func Headlines(db database.Database, f database.Filter, w io.Writer) error {

	file := xlsx.NewFile()

	g, err := db.GroupByFilter(f)
	if err != nil {
		return err
	}
	sheet, err := file.AddSheet("Broadsheet")
	if err != nil {
		return err
	}
	headlineBroadsheet(sheet, g)

	sheet, err = file.AddSheet("Summary")
	if err != nil {
		return err
	}
	exportInfoSheet(sheet, db, f)

	err = file.Write(w)
	if err != nil {
		return err
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

func headlineBroadsheet(sheet *xlsx.Sheet, g group.Group) error {

	// Subject mapped by ebacc
	subjMap := map[subject]struct{}{}
	for _, s := range g.Students {
		for _, r := range s.Results {
			header := subject{r.Subj, r.EBacc}
			subjMap[header] = struct{}{}
		}
	}

	subjects := subjList{}
	for s, _ := range subjMap {
		subjects = append(subjects, s)
	}

	sort.Sort(subjects)

	colours := map[string]string{
		"En": "FF8DB4E2",
		"El": "FF8DB4E2",
		"M":  "FF0066CC",
		"S":  "FF73FF47",
		"H":  "FF73FF47",
		"L":  "FF73FF47",
		"":   "FFFFCC00",
	}

	row := sheet.AddRow()
	row.SetHeightCM(4.5)
	cell := row.AddCell()
	cell.Value = "Name"
	for _, s := range subjects {
		cell := row.AddCell()
		cell.Value = s.Name
		style := xlsx.NewStyle()
		style.Font = defaultFont
		style.Fill.PatternType = "solid"
		style.Fill.FgColor = colours[s.EBacc]
		style.ApplyFill = true
		style.Alignment.TextRotation = 90
		style.Alignment.Horizontal = "center"
		style.ApplyAlignment = true
		cell.SetStyle(style)
	}

	firstCol := 1
	err := sheet.SetColWidth(firstCol, firstCol+len(subjects), 2.6)
	if err != nil {
		return err
	}

	return nil
}
