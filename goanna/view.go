package goanna

import (
	"github.com/99designs/goodies/sprockets"
	"html/template"
)

type Views struct {
	*template.Template
	AssetPipeline *sprockets.ViewHelper
}

func NewView(viewHelper *sprockets.ViewHelper) Views {
	v := Views{nil, viewHelper}
	if v.Template == nil {
		v.Template = template.New("Template")
	}

	v.Template.Funcs(template.FuncMap{
		"urlFor":     UrlFor,
		"stylesheet": v.AssetPipeline.StylesheetLinkTag,
		"javascript": v.AssetPipeline.JavascriptTag,
		"assetUrl":   v.AssetPipeline.GetAssetUrl,
	})

	return v
}
