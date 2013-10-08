package goanna

import (
	"github.com/gorilla/mux"
	"log"
	"net/url"
)

// Router is a global router used for routing requests
var Router *mux.Router

// UrlBase is the base url used for creating absolute urls
var UrlBase url.URL

func init() {
	ClearRouter()
	UrlBase.Scheme = "http"
	UrlBase.Host = "localhost"
}

// ClearRouter clears the global router
func ClearRouter() {
	Router = mux.NewRouter()
}

// UrlFor returns the relative URL for the route name
func UrlFor(name string, args ...string) *url.URL {
	r := Router.Get(name)
	if r == nil {
		log.Panicf("UrlFor: No such route '%s'\n", name)
	}
	u, err := r.URL(args...)
	if err != nil {
		log.Panicln("UrlFor: " + err.Error())
	}
	return u
}

// UrlFor returns the avsolute URL for the route name.
// UrlBase is used for the URL Host and Scheme
func AbsoluteUrlFor(name string, args ...string) *url.URL {
	u := UrlFor(name, args...)
	u.Host = UrlBase.Host
	u.Scheme = UrlBase.Scheme
	return u
}
