package template

import (
	"bytes"
	"strings"
	"text/template"
)

func Exec(name string, raw string, payload any) (*string, error) {
	tpl, err := template.
		New(name).
		Funcs(template.FuncMap{
			"join": func(sep string, s ...string) string {
				return strings.Join(s, sep)
			},
		}).
		Parse(raw)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	err = tpl.
		ExecuteTemplate(&buffer, name, payload)
	if err != nil {
		return nil, err
	}

	templated := buffer.String()
	return &templated, nil
}
