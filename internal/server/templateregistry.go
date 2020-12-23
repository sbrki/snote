package server

import (
	"io"
	"text/template"

	"github.com/labstack/echo"
)

type TemplateRegistry struct {
	templates *template.Template
}

func NewTemplateRegistry(templatesPath string) *TemplateRegistry {
	tr := new(TemplateRegistry)
	tr.templates = template.Must(template.ParseGlob(templatesPath))
	return tr
}

func (tr *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return tr.templates.ExecuteTemplate(w, name, data)
}
