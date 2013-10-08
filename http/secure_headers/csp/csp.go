// package csp implements a content-security-policy header generator
package csp

import (
	"strings"
)

// See CSP standard at http://www.w3.org/TR/CSP/
const (
	// Allows content to be loaded from the current domain
	SourceSelf = "'self'"
	// Prevents any content of the specified type loading
	SourceNone = "'none'"
	// Disables the main protection offered by CSP
	SourceUnsafeInline = "'unsafe-inline'"
)

// Opts configures a Content-Security-Policy header
type Opts struct {
	ReportUri string // A relative path to POST CSP violations to

	// Sources to allow content loading from.
	DefaultSrc []string
	ScriptSrc  []string
	ConnectSrc []string
	FrameSrc   []string
	FontSrc    []string
	ImgSrc     []string
	MediaSrc   []string
	ObjectSrc  []string
	StyleSrc   []string
}

// Formats for rendering as an http header.
// E.G. default-src 'self' ; script-src 'self' https://apis.google.com
func (o Opts) Header() string {
	h := make([]string, 0)

	var fields = []struct {
		desc  string
		sites []string
	}{
		{"default-src", o.DefaultSrc},
		{"connect-src", o.ConnectSrc},
		{"frame-src", o.FrameSrc},
		{"font-src", o.FontSrc},
		{"img-src", o.ImgSrc},
		{"media-src", o.MediaSrc},
		{"object-src", o.ObjectSrc},
		{"style-src", o.StyleSrc},
		{"script-src", o.ScriptSrc},
	}
	for _, f := range fields {
		if f.sites != nil {
			h = append(h, f.desc+" "+strings.Join(f.sites, " "))
		}
	}
	if o.ReportUri != "" {
		h = append(h, "report-uri "+o.ReportUri)
	}
	return strings.Join(h, " ; ")
}
