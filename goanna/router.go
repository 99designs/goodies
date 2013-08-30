package goanna

import (
	"github.com/gorilla/mux"
	"log"
	"net/url"
)

var Router *mux.Router
var UrlBase url.URL

func init() {
	ClearRouter()
	UrlBase.Scheme = "http"
	UrlBase.Host = "localhost"
}

func ClearRouter() {
	// For cleanup between test runs
	Router = mux.NewRouter()
}

func UrlFor(name string, args ...string) *url.URL {
	r := Router.Get(name)
	if r == nil {
		log.Panicf("UrlFor: No such route '%s'\n", name)
	}
	url, err := r.URL(args...)
	if err != nil {
		log.Panicln("UrlFor: " + err.Error())
	}
	return url
}

func AbsoluteUrlFor(name string, args ...string) *url.URL {
	u := UrlFor(name, args...)
	u.Host = UrlBase.Host
	u.Scheme = UrlBase.Scheme
	return u
}
