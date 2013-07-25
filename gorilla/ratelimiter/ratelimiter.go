// Package ratelimiter decorates an gorilla.Mux and prevents too many requests being made
package ratelimiter

import (
	rl "github.com/99designs/goodies/http/ratelimiter"
	"github.com/gorilla/mux"
)

type RateLimitedHandler struct {
	*mux.Router
	rl.RateLimitedHandler
}

func Decorate(delegate *mux.Router, rateLimiter rl.RateLimiter) RateLimitedHandler {
	return RateLimitedHandler{delegate, rl.Decorate(delegate, rateLimiter)}
}
