package queue

import (
	"bytes"
	"html/template"
)

type TemplateFile string

func (tf TemplateFile) String() string {
	return string(tf)
}

const (
	templateIndex TemplateFile = "templates/notification/email/index.html"
)

type TemplateGenerator struct{}

func NewTemplateGenerator() *TemplateGenerator {
	return &TemplateGenerator{}
}

func (g *TemplateGenerator) ToString(templateFile TemplateFile, data any) (string, error) {
	tmpl, err := template.ParseFiles(templateIndex.String(), templateFile.String())
	if err != nil {
		return "", err
	}

	var bodyBuffer bytes.Buffer
	if err := tmpl.Execute(&bodyBuffer, data); err != nil {
		return "", err
	}

	return bodyBuffer.String(), nil
}
