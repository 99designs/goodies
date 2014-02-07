package assets

import (
	"html/template"
	"net/http"
)

type AssetPipeline interface {
	GetAssetUrl(name string) (string, error)
	GetAssetContents(name string) ([]byte, error)
}

type AssetHelper struct {
	AssetPipeline
}

func (vh *AssetHelper) StylesheetLinkTag(name string) (template.HTML, error) {
	url, err := vh.AssetUrl(name)
	return template.HTML(`<link href="` + template.HTMLEscaper(url) + `" rel="stylesheet" type="text/css">`), err
}

func (vh *AssetHelper) InlineStylesheet(name string) (template.HTML, error) {
	content, err := vh.AssetPipeline.GetAssetContents(name)
	return template.HTML(`<style>` + string(content) + `</style>`), err
}

func (vh *AssetHelper) JavascriptTag(name string) (template.HTML, error) {
	url, err := vh.AssetUrl(name)
	return template.HTML(`<script src="` + template.HTMLEscaper(url) + `" type="text/javascript"></script>`), err
}

func (vh *AssetHelper) AssetUrl(name string) (template.HTMLAttr, error) {
	url, err := vh.AssetPipeline.GetAssetUrl(name)
	return template.HTMLAttr(url), err
}

type AssetHandler struct {
	AssetPipeline
}

func (vh AssetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := vh.GetAssetContents(r.URL.Path[1:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.Write(c)
	}
}
