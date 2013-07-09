// Package http rate decorates an http.Handler and prevents too many requests being made
package http

import (
	"net/http"
)

const StatusTooManyRequests = 429

type RateLimitedHandler struct {
	http.Handler
	rateLimiter RateLimiter
}

type RateLimiter func(*http.Request) bool

func Decorate(delegate http.Handler, rateLimiter RateLimiter) RateLimitedHandler {
	return RateLimitedHandler{delegate, rateLimiter}
}

func (lh RateLimitedHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if lh.rateLimiter(r) {
		lh.Handler.ServeHTTP(rw, r)
	} else {
		http.Error(rw, "Too many requests", StatusTooManyRequests)
	}
}
