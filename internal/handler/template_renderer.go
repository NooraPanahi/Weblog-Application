package handler

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func NewTemplateRenderer() *TemplateRenderer {
	templatee := template.Must(template.ParseGlob("templates/**/*.html"))

	return &TemplateRenderer{templates: templatee}
}

func (r* TemplateRenderer) Render (w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name,data)
}