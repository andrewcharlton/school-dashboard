package templates

import "html/template"

// Parse the contents of the template files and create an html.Template
func Parse() (*template.Template, error) {

	var t *template.Template
	for name, contents := range allTemplates {
		var tmpl *template.Template

		if t == nil {
			t = template.New(name)
		}

		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}

		_, err := tmpl.Parse(contents)
		if err != nil {
			return t, err
		}
	}

	return t, nil
}
