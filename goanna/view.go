package goanna

import (
	"github.com/99designs/goodies/sprockets"
	"html/template"
)

type Views struct {
	*template.Template
	ViewHelper *sprockets.ViewHelper
}

func NewView(viewHelper *sprockets.ViewHelper) Views {
	tmpl := template.New("Template")

	return Views{tmpl, viewHelper}
}
