package sprockets

import (
	"html/template"
)

func NullSprockets() ViewHelper {
	return ViewHelper{nullSprockets{}}
}

type nullSprockets struct{}

func (s nullSprockets) GetAssetUrl(name string) (template.HTMLAttr, error) {
	return "", nil
}

func (s nullSprockets) GetAssetContents(name string) ([]byte, error) {
	return make([]byte, 0), nil
}
