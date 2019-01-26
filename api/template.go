package api

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

type Template struct {
}

func (tpl *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	t, err := template.New("").ParseFiles("views/"+name+".html", "views/base.html")
	if err != nil {
		return err
	}

	if obj := t.Lookup("js"); obj == nil {
		t, err = t.Parse(`{{define "js"}}{{end}}`)
		if err != nil {
			return err
		}
	}

	if obj := t.Lookup("css"); obj == nil {
		t, err = t.Parse(`{{define "css"}}{{end}}`)
		if err != nil {
			return err
		}
	}

	return t.ExecuteTemplate(w, "base", data)
}
