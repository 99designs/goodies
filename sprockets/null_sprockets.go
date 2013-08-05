package sprockets

import (
	"html/template"
)

func NullSprockets() ViewHelper {
	return ViewHelper{nullSprockets{}}
}

type nullSprockets struct{}

func (s nullSprockets) GetAssetUrl(name string) (template.HTMLAttr, error) {
	panic("Null sprockets in use, cannot get asset urls")
}

func (s nullSprockets) GetAssetContents(name string) ([]byte, error) {
	panic("Null sprockets in use, cannot get asset contents")
}
