package csp

import (
	"fmt"
	_ "testing"
)

func ExampleOpts() {
	h := Opts{
		DefaultSrc: []string{SourceNone},
		ImgSrc:     []string{SourceSelf, "https://example.org"},
		StyleSrc:   []string{SourceSelf},
		ScriptSrc:  []string{SourceSelf, SourceUnsafeInline},
		ReportUri:  "/csp_report",
	}
	fmt.Println(h.Header())
	// Output: default-src 'none' ; img-src 'self' https://example.org ; style-src 'self' ; script-src 'self' 'unsafe-inline' ; report-uri /csp_report
}
