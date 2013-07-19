// Package http rate decorates an gorilla.Mux and prevents too many requests being made
package gorilla

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const StatusTooManyRequests = 429

type RateLimitedHandler struct {
	*mux.Router
	rateLimiter RateLimiter
}

type RateLimiter func(*http.Request) (ok bool, retryAfter time.Duration)

func Decorate(delegate *mux.Router, rateLimiter RateLimiter) RateLimitedHandler {
	return RateLimitedHandler{delegate, rateLimiter}
}

func (lh RateLimitedHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if ok, retryAfter := lh.rateLimiter(r); ok {
		lh.Router.ServeHTTP(rw, r)
	} else {
		rw.Header().Set("Retry-After", fmt.Sprint(int64(retryAfter.Seconds())));
		http.Error(rw, "Too many requests", StatusTooManyRequests)
	}
}
