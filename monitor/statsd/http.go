package statsd

import (
	"net/http"
	"time"
)

type StatsdLogHandler struct {
	handler http.Handler
	statter *Statsd
	bucket  string
}

// ServeHTTP records the time for a request and sends the result to statsd
func (sh *StatsdLogHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	requestStartTime := time.Now()
	sh.handler.ServeHTTP(w, req)
	requestDuration := time.Now().Sub(requestStartTime)

	sh.statter.Counter(1.0, sh.bucket+".requests", 1)
	sh.statter.Timing(1.0, sh.bucket+".responsetime", requestDuration)
}

func Decorate(statter *Statsd, bucket string, h http.Handler) http.Handler {
	return &StatsdLogHandler{
		handler: h,
		statter: statter,
		bucket:  bucket,
	}
}
