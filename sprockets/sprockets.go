/*
Package sprockets integrates the power of sprockets ( http://getsprockets.org ) with your go program.

To use in development:
 * Assets go in '/assets'
 * Install the 'rack' and 'sprockets' rubygems
 * run `rackup src/github.com/99designs/goodies/sprockets/config.ru`
 * In your app, call 'NewSprocketsServer("http://localhost:8012", "assets") to get a ViewHelper'
 * Use the ViewHelper to generate links to your assets.
*/
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

func (vh *ViewHelper) JavascriptTag(name string) (template.HTML, error) {
	var attr template.HTMLAttr
	url, err := vh.AssetPipeline.GetAssetUrl(name)
	if err != nil {
		return template.HTML(""), err
	}
	attr = template.HTMLAttr(url)
	return template.HTML(`<script src="` + template.HTMLEscaper(attr) + `" type="text/javascript"></script>`), nil
}
