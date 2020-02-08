package main

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"text/template"
)

func renderToBuffer(config interface{}, tpl string, name string, buffer bytes.Buffer) bytes.Buffer {
	tmpl := template.Must(template.New(name).Funcs(template.FuncMap{
		"determineVariableType": func(variable string) string {
			return "string"
		},
	}).Parse(tpl))
	err := tmpl.Execute(&buffer, config)

	if err != nil {
		logrus.Fatal("Failed writing to Buffer: ", err)
	}

	return buffer
}
