// Decorates an http.Handler with DOS protection from auth/ratelimiter
package http

import (
	"net/http"
	"strings"
)

const StatusTooManyRequests = 429

type RateLimitedHandler struct {
	http.Handler
	rateLimiter RateLimiter
}

type RateLimiter func(*http.Request) bool

func Decorate(delegate http.Handler, rateLimiter RateLimiter) RateLimitedHandler {
	if rateLimiter == nil {
		rateLimiter = DefaultRateLimiter
	}
	return RateLimitedHandler{delegate, rateLimiter}
}

func (lh RateLimitedHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if lh.rateLimiter(r) {
		lh.Handler.ServeHTTP(rw, r)
	} else {
		http.Error(rw, "Too many requests", StatusTooManyRequests)
	}
}
