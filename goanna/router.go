package goanna

import (
	"github.com/gorilla/mux"
	"log"
	"net/url"
)

var Router *mux.Router

func init() {
	ClearRouter()
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
