// Package cachecontrol provides an HTTP decorator which adds cache-control headers to outgoing responses
package cachecontrol

import (
	"fmt"
	"net/http"
)

type cacheDecorator struct {
	delegate http.Handler
	maxage   uint
}

func CacheControl(maxage uint, delegate http.Handler) http.Handler {
	return &cacheDecorator{delegate, maxage}
}

func (f *cacheDecorator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", f.maxage))
	f.delegate.ServeHTTP(w, r)
}
