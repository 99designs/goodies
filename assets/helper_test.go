package assets

import (
	"html/template"
	"testing"
	. "launchpad.net/gocheck"
)

type ST struct{}

var _ = Suite(&ST{})

func Test(t *testing.T) { TestingT(t) }

func (s *ST) TestAssetHelperMethods(c *C) {
	link := StylesheetTag("&Foo")
	c.Check(link, Equals, template.HTML(`<link href="&amp;Foo" rel="stylesheet" type="text/css">`))

	js := JavascriptTag("&Foo")
	c.Check(js, Equals, template.HTML(`<script src="&amp;Foo" type="text/javascript"></script>`))
}
