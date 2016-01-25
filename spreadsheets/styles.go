package spreadsheets

import "github.com/tealeg/xlsx"

// Fonts
var fonts = map[string]xlsx.Font{
	"Default": xlsx.Font{11, "Calibri", 2, 0, "FF000000", false, false, false},
	"Bold":    xlsx.Font{11, "Calibri", 2, 0, "FF000000", true, false, false},
	"Title":   xlsx.Font{18, "Calibri", 2, 0, "FF000000", false, false, false},
}

// Fills
var fills = map[string]xlsx.Fill{
	"None":      xlsx.Fill{},
	"LightGrey": xlsx.Fill{FgColor: "FFD9D9D9", PatternType: "solid"},
	"LightBlue": xlsx.Fill{FgColor: "FF8DB4E2", PatternType: "solid"},
	"DarkBlue":  xlsx.Fill{FgColor: "FF0066CC", PatternType: "solid"},
	"Green":     xlsx.Fill{FgColor: "FF73FF47", PatternType: "solid"},
	"Yellow":    xlsx.Fill{FgColor: "FFFFCC00", PatternType: "solid"},
	"Orange":    xlsx.Fill{FgColor: "FFE6822D", PatternType: "solid"},
	"Purple":    xlsx.Fill{FgColor: "FFF660D9", PatternType: "solid"},
	"Red":       xlsx.Fill{FgColor: "FFF7774F", PatternType: "solid"},
}

// Borders
var borders = map[string]xlsx.Border{
	"None":   xlsx.Border{},
	"Bottom": xlsx.Border{Bottom: "thin"},
	"All":    xlsx.Border{Left: "thin", Right: "thin", Top: "thin", Bottom: "thin"},
}

// Alignments
var alignments = map[string]xlsx.Alignment{
	"Left":     xlsx.Alignment{},
	"Center":   xlsx.Alignment{Horizontal: "center"},
	"Vertical": xlsx.Alignment{Horizontal: "center", TextRotation: 90},
}

func newStyle(font, fill, border, align string) *xlsx.Style {

	return &xlsx.Style{borders[border], fills[fill], fonts[font], true, true, true,
		true, alignments[align]}
}

var left = &xlsx.Style{borders["None"], fills["None"], fonts["Default"], false,
	false, true, false, alignments["Left"]}

var center = &xlsx.Style{borders["None"], fills["None"], fonts["Default"], false,
	false, true, true, alignments["Center"]}
