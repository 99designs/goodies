// Package secure_headers decorates an http.Handler and
// sets several key security headers
package secure_headers

import (
	"github.com/99designs/goodies/http/secure_headers/csp"
	"net/http"
)

// Settings configures the headers a secure handler will add to a ResponseWriter
type Settings struct {
	CspOpts    csp.Opts // Content-Security-Policy
	ReportOpts csp.Opts // Content-Security-Policy-Report-Only

	// Where can this site be embedded as an iframe
	FrameOptions string
	// Should user agents default to SSL
	StrictTransportSecurity string
	// Should IE guess mime types
	ContentTypeOptions string
	// Should IE run code that 'looks like' an XSS
	XSSProtection string
	// Specify which cross-domain policies flash can load
	PermittedCrossDomainPolicies string
}

// Sane/safe defaults for the secure headers decorator.
// Content-Security-Policy is disabled by default
// as it is very restrictive.
var DefaultSettings = Settings{
	CspOpts:                      csp.Opts{},
	ReportOpts:                   csp.Opts{},
	FrameOptions:                 "SAMEORIGIN",
	StrictTransportSecurity:      "max-age=31536000; includeSubDomains",
	ContentTypeOptions:           "nosniff",
	XSSProtection:                "1; mode=block",
	PermittedCrossDomainPolicies: "master-only",
}

func Decorate(settings Settings, delegate http.Handler) http.Handler {
	return secureHandler{
		Delegate: delegate,
		Settings: settings,
	}
}

func (s Settings) setHeaders(rw http.ResponseWriter) {
	for _, h := range []struct {
		key   string
		value string
	}{
		{"X-Frame-Options", s.FrameOptions},
		{"Strict-Transport-Security", s.StrictTransportSecurity},
		{"X-Content-Type-Options", s.ContentTypeOptions},
		{"X-XSS-Protection", s.XSSProtection},
		{"X-Permitted-Cross-Domain-Policies", s.PermittedCrossDomainPolicies},
		{"Content-Security-Policy", s.CspOpts.Header()},
		{"Content-Security-Policy-Report-Only", s.ReportOpts.Header()},
	} {
		rw.Header().Add(h.key, h.value)
	}
}

// Type secureHandler adds security headers to responses before passing them through to the delegate handler.
type secureHandler struct {
	Delegate http.Handler
	Settings
}

func (lh secureHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	lh.setHeaders(rw)
	lh.Delegate.ServeHTTP(rw, r)
}
