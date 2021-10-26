package utils

import (
	"bytes"
	"html/template"
)

// RenderString
func RenderString(tpl string, vars interface{}) (string, error) {
	t, err := template.New("tpl").Parse(tpl)

	if err != nil {
		return "", err
	}
	return process(t, vars)
}

// process
func process(t *template.Template, vars interface{}) (string, error) {
	var tmplBytes bytes.Buffer

	err := t.Execute(&tmplBytes, vars)
	if err != nil {
		return "", err
	}
	return tmplBytes.String(), nil
}
