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
	url, err := vh.asset_url(name)
	return template.HTML(`<link href="` + template.HTMLEscaper(url) + `" rel="stylesheet">`), err
}

func (vh *ViewHelper) InlineStylesheet(name string) (template.HTML, error) {
	content, err := vh.AssetPipeline.GetAssetContents(name)
	return template.HTML(`<style>` + string(content) + `</style>`), err
}

func (vh *ViewHelper) JavascriptTag(name string) (template.HTML, error) {
	url, err := vh.asset_url(name)
	return template.HTML(`<script src="` + template.HTMLEscaper(url) + `" type="text/javascript"></script>`), err
}

func (vh *ViewHelper) asset_url(name string) (template.HTMLAttr, error) {
	url, err := vh.AssetPipeline.GetAssetUrl(name)
	return template.HTMLAttr(url), err
}
