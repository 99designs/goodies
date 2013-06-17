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
	v := Views{nil, viewHelper}
	v.Init()
	return v
}

func (v *Views) Init() {
	if v.Template == nil {
		v.Template = template.New("Template")
	}

	v.Template.Funcs(template.FuncMap{
		"urlFor":     UrlFor,
		"stylesheet": v.ViewHelper.StylesheetLinkTag,
		"javascript": v.ViewHelper.JavascriptTag,
	})
}
