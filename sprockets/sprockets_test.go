package sprockets

import (
	"html/template"
	. "launchpad.net/gocheck"
	"testing"
)

type TestAssetPipeline struct {
}

func (t TestAssetPipeline) GetAssetUrl(name string) (template.HTMLAttr, error) {
	return template.HTMLAttr(name), nil
}

func (t TestAssetPipeline) GetAssetContents(name string) ([]byte, error) {
	return []byte("Contents of " + name), nil
}

type ST struct{}

var _ = Suite(&ST{})

func Test(t *testing.T) { TestingT(t) }

func (s *ST) TestAssetHelperMethods(c *C) {
	vh := ViewHelper{TestAssetPipeline{}}
	inline, err := vh.InlineStylesheet("&Foo")
	c.Check(inline, Equals, template.HTML(`<style>Contents of &Foo</style>`))
	c.Check(err, Equals, nil)

	link, err := vh.StylesheetLinkTag("&Foo")
	c.Check(err, Equals, nil)
	c.Check(link, Equals, template.HTML(`<link href="&amp;Foo" rel="stylesheet">`))

	js, err := vh.JavascriptTag("&Foo")
	c.Check(err, Equals, nil)
	c.Check(js, Equals, template.HTML(`<script src="&amp;Foo" type="text/javascript"></script>`))
}
