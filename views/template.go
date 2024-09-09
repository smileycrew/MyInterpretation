package views

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

type Template struct {
	templates *template.Template // *template.Template is a type which is used to represent one of more parsed templates
}

func GetTemplates() *Template {
	tmpl := template.Must(template.ParseGlob("*templates/layouts/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("*templates/pages/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("*templates/partials/*.html"))

	return &Template{
		templates: tmpl,
	}
}

// function called Render (based on Template) takes io writer, name, data, context from echo and returns an error
func (template *Template) Render(writer io.Writer, name string, data interface{}, context echo.Context) error {
	// executes template by writing to writer the template and the data
	return template.templates.ExecuteTemplate(writer, name, data)
}
