// Package log decorates any http handler with a
// Common Log Format logger
package log

import (
	"log"
	"net/http"
	"os"
	"time"
)

type commonLogHandler struct {
	handler http.Handler
	logger  *log.Logger
}

// CommonLogHandler returns a handler that serves HTTP requests
// If a logger is not provided, stdout will be used
func CommonLogHandler(logger *log.Logger, h http.Handler) http.Handler {
	if logger == nil {
		logger = log.New(os.Stdout, "", 0)
	}
	return &commonLogHandler{
		handler: h,
		logger:  logger,
	}
}

// ServeHTTP logs the request and response data to Common Log Format
func (lh *commonLogHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// grab the data we need before passing it to ServeHTTP
	startTime := time.Now()
	reqRemoteAddr := req.RemoteAddr
	reqURLUserUsername := "-"
	if req.URL.User != nil {
		if name := req.URL.User.Username(); name != "" {
			reqURLUserUsername = name
		}
	}
	reqMethod := req.Method
	reqRequestURI := req.RequestURI
	reqProto := req.Proto

	// decorate the writer so we can capture the status and size
	loggedWriter := &responseLogger{w: w}
	lh.handler.ServeHTTP(loggedWriter, req)

	// Common Log Format
	lh.logger.Printf("%s %s - [%s] \"%s %s %s\" %d %d",
		reqRemoteAddr,
		reqURLUserUsername,
		startTime.Format("02/Jan/2006:15:04:05 -0700"),
		reqMethod,
		reqRequestURI,
		reqProto,
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
