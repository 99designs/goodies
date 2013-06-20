// Package log decorates any http handler with a
// Common Log Format logger
package log

import (
	"log"
	"net/http"
	"os"
	"time"
)

type CommonLogHandler struct {
	handler http.Handler
	logger  *log.Logger
}

func DefaultCommonLogHandler(h http.Handler) http.Handler {
	return &CommonLogHandler{
		handler: h,
		logger:  log.New(os.Stderr, "", 0),
	}
}

func NewCommonLogHandler(logger *log.Logger, h http.Handler) http.Handler {
	return &CommonLogHandler{
		handler: h,
		logger:  logger,
	}
}

func (lh *CommonLogHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()

	loggedWriter := &responseLogger{w: w}
	lh.handler.ServeHTTP(loggedWriter, req)

	// Common Log Format
	username := "-"
	if req.URL.User != nil {
		if name := req.URL.User.Username(); name != "" {
			username = name
		}
	}

	lh.logger.Printf("%s %s - [%s] \"%s %s %s\" %d %d",
		req.RemoteAddr,
		username,
		startTime.Format("02/Jan/2006:15:04:05 -0700"),
		req.Method,
		req.RequestURI,
		req.Proto,
		loggedWriter.status,
		loggedWriter.size,
	)
}

type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	if l.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		l.status = http.StatusOK
	}
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}
