package templates

import (
	"html/template"
	"strconv"
)

var funcMap = template.FuncMap{"Percent": percentage}

func percentage(x float64, dp int) string {

	return strconv.FormatFloat(100.0*x, "f", dp, 64)

}
