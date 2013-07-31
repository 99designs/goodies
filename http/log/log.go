// Package log decorates any http handler with a
// Common Log Format logger
package log

import (
	"fmt"
	"github.com/99designs/goodies/http/log/response"
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
	// since ServeHTTP could modify that data
	data := NewCommonLogData(*req)

	// decorate the writer so we can capture the status and size
	loggedWriter := response.LogResponseMetadata(w)
	lh.handler.ServeHTTP(loggedWriter, req)

	data.ResponseStatus = loggedWriter.Status
	data.ResponseSize = loggedWriter.Size
	lh.logger.Println(data.String())
}

type CommonLogData struct {
	StartTime      time.Time
	RemoteAddr     string
	Username       string
	Method         string
	Uri            string
	Proto          string
	ResponseStatus int
	ResponseSize   int
}

func NewCommonLogData(req http.Request) CommonLogData {
	return CommonLogData{
		StartTime:  time.Now(),
		RemoteAddr: req.RemoteAddr,
		Username:   usernameFromReq(req),
		Method:     req.Method,
		Uri:        req.RequestURI,
		Proto:      req.Proto,
	}
}

func (this CommonLogData) String() string {
	return fmt.Sprintf("%s %s - [%s] \"%s %s %s\" %d %d",
		this.RemoteAddr,
		this.Username,
		this.StartTime.Format("02/Jan/2006:15:04:05 -0700"),
		this.Method,
		this.Uri,
		this.Proto,
		this.ResponseStatus,
		this.ResponseSize,
	)
}

// Extract username
func usernameFromReq(req http.Request) string {
	if req.URL.User != nil {
		if name := req.URL.User.Username(); name != "" {
			return name
		}
	}
	return "-"
}
