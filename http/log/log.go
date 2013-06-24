// Package log decorates any http handler with a
// Common Log Format logger
package log

import (
	"log"
	"net/http"
	"os"
	"time"
)

const (
	UseXForwardedAddress = 1 << iota
)

type CommonLogHandler struct {
	handler http.Handler
	logger  *log.Logger
	flags   int
}

// flags_arr is so that no code is broken by an API change
func DefaultCommonLogHandler(h http.Handler, flags_arr ...int) http.Handler {
	flags := 0
	if len(flags_arr) > 0 {
		flags = flags_arr[0]
	}
	return &CommonLogHandler{
		handler: h,
		logger:  log.New(os.Stderr, "", log.LstdFlags),
		flags:   flags,
	}
}

func NewCommonLogHandler(logger *log.Logger, h http.Handler, flags_arr ...int) http.Handler {
	flags := 0
	if len(flags_arr) > 0 {
		flags = flags_arr[0]
	}
	return &CommonLogHandler{
		handler: h,
		logger:  logger,
		flags:   flags,
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
	addr := req.RemoteAddr
	if (lh.flags & UseXForwardedAddress) != 0 {
		ip := req.Header.Get("X-Forwarded-For")
		port := req.Header.Get("X-Forwarded-Port")
		if ip != "" && port != "" {
			addr = ip + ":" + port
		}
	}
	lh.logger.Printf("%s %s - [%s] \"%s %s %s\" %d %d",
		addr,
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
