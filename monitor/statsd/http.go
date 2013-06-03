package statsd

import (
	"fmt"
	"net/http"
	"time"
)

type StatsdLogHandler struct {
	handler http.Handler
	statter *Statsd
	prefix  string
}

// ServeHTTP records the time for a request and sends the result to statsd
func (sh *StatsdLogHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	requestStartTime := time.Now()
	sh.handler.ServeHTTP(w, req)
	requestDuration := time.Now().Sub(requestStartTime)

	description := fmt.Sprintf("%s.%s.%s.", sh.prefix, req.Method, req.URL.Path)
	sh.statter.Counter(1.0, description+"requests", 1)
	sh.statter.Timing(1.0, description+"responsetime", requestDuration)
}

func Decorate(statter *Statsd, prefix string, h http.Handler) http.Handler {
	return &StatsdLogHandler{
		handler: h,
		statter: statter,
		prefix:  prefix,
	}
}
