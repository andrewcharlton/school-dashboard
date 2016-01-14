package spreadsheets

import "github.com/tealeg/xlsx"

// Fills
var noFill = xlsx.Fill{}

// Fonts
var boldFont = xlsx.Font{11, "Calibri", 2, 0, "FF000000", true, false, false}
var defaultFont = xlsx.Font{11, "Calibri", 2, 0, "FF000000", false, false, false}
var titleFont = xlsx.Font{18, "Calibri", 2, 0, "FF000000", false, false, false}

// Borders
var noBorder = xlsx.Border{}
var bottomBorder = xlsx.Border{Bottom: "thin"}
var allBorders = xlsx.Border{Left: "thin", Right: "thin", Top: "thin", Bottom: "thin"}

// Alignments
var left = xlsx.Alignment{}
var center = xlsx.Alignment{Horizontal: "center"}
var vert = xlsx.Alignment{TextRotation: 90}

// Styles
var defaultStyle = &xlsx.Style{noBorder, noFill, defaultFont, false, false, true, false, left}
var centered = &xlsx.Style{noBorder, noFill, defaultFont, false, false, true, true, center}
var bold = &xlsx.Style{noBorder, noFill, boldFont, false, false, true, false, left}
var title = &xlsx.Style{noBorder, noFill, titleFont, false, false, true, false, left}
var header = &xlsx.Style{bottomBorder, noFill, boldFont, true, false, true, false, left}
var centerHeader = &xlsx.Style{bottomBorder, noFill, boldFont, true, false, true, true, center}
var vertHeader = &xlsx.Style{bottomBorder, noFill, boldFont, true, false, true, true, vert}
var gridStyle = &xlsx.Style{allBorders, noFill, defaultFont, true, false, true, true, center}

// Used to duplicate a style, without effecting original
func copyStyle(s *xlsx.Style) *xlsx.Style {

	style := xlsx.Style{s.Border, s.Fill, s.Font, s.ApplyBorder, s.ApplyFill, s.ApplyFont,
		s.ApplyAlignment, s.Alignment}
	return &style
}
