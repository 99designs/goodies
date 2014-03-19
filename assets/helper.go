package assets

import (
	"fmt"
	"html/template"
)

func StylesheetTag(url string) template.HTML {
	return template.HTML(
		fmt.Sprintf(`<link href="%s" rel="stylesheet" type="text/css">`, template.HTMLEscaper(url)))
}

func JavascriptTag(url string) template.HTML {
	return template.HTML(
		fmt.Sprintf(`<script src="%s" type="text/javascript"></script>`, template.HTMLEscaper(url)))
}
