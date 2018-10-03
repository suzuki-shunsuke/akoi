package util

import (
	"bytes"
	"text/template"
)

// RenderTpl renders a template.
func RenderTpl(tpl *template.Template, data interface{}) (string, error) {
	buf := bytes.Buffer{}
	err := tpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
