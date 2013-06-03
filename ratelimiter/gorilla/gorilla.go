package gorilla

import (
	"github.com/gorilla/mux"
	"net/http"
)

const StatusTooManyRequests = 429

type RateLimitedHandler struct {
	*mux.Router
	rateLimiter RateLimiter
}

type RateLimiter func(*http.Request) bool

func Decorate(delegate *mux.Router, rateLimiter RateLimiter) RateLimitedHandler {
	return RateLimitedHandler{delegate, rateLimiter}
}

func (lh RateLimitedHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if lh.rateLimiter(r) {
		lh.Router.ServeHTTP(rw, r)
	} else {
		http.Error(rw, "Too many requests", StatusTooManyRequests)
	}
}
