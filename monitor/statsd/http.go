package statsd

import (
	"fmt"
	"github.com/99designs/goodies/goanna"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type StatsdLogHandler struct {
	handler          http.Handler
	statter          *Statsd
	prefix           string
	RequestFormatter func(req *http.Request) string
}

// ServeHTTP records the time for a request and sends the result to statsd
func (sh *StatsdLogHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	requestStartTime := time.Now()
	sh.handler.ServeHTTP(w, req)
	requestDuration := time.Now().Sub(requestStartTime)

	description := sh.RequestFormatter(req)
	if description != "" {
		sh.statter.Counter(1.0, sh.prefix+"."+description+".requests", 1)
		sh.statter.Timing(1.0, sh.prefix+"."+description+".responsetime", requestDuration)
	}
}

func defaultGetRequestString(req *http.Request) string {
	var match mux.RouteMatch
	if goanna.Router.Match(req, &match) {
		name := match.Route.GetName()
		if name != "" {
			return fmt.Sprintf("%s.%s", req.Method, name)
		}
	}
	return ""
}

func Decorate(statter *Statsd, prefix string, h http.Handler) http.Handler {
	return &StatsdLogHandler{
		handler:          h,
		statter:          statter,
		prefix:           prefix,
		RequestFormatter: defaultGetRequestString,
	}
}
