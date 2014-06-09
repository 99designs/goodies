package http

import (
	"net/http"
	"regexp"
)

// RegexPath returns a handler that serves HTTP requests
// with the first match of the the given regex from the request URL's Path
// and invoking the handler h. RegexPath handles a
// request for a path that doesn't match by
// replying with an HTTP 404 not found error.
func RegexPath(reg *regexp.Regexp, h http.Handler) http.Handler {
	if reg == nil {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		matches := reg.FindStringSubmatch(r.URL.Path)
		if len(matches) > 1 {
			r.URL.Path = matches[1]
			h.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}
