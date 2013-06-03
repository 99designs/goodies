package sprockets

import (
	"html/template"
)

type AssetPipeline interface {
	GetAssetUrl(name string) (template.HTMLAttr, error)
	GetAssetContents(name string) ([]byte, error)
}

type ViewHelper struct {
	AssetPipeline
}

func (vh *ViewHelper) StylesheetLinkTag(name string) (template.HTML, error) {
	var attr template.HTMLAttr
	url, err := vh.AssetPipeline.GetAssetUrl(name)
	if err != nil {
		return template.HTML(""), err
	}
	attr = template.HTMLAttr(url)
	return template.HTML(`<link href="` + template.HTMLEscaper(attr) + `" rel="stylesheet">`), nil
}
