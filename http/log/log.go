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
	handler    http.Handler
	logger     *log.Logger
	newLogData func() CommonLogData
}

// CommonLogHandler returns a handler that serves HTTP requests
// If a logger is not provided, stdout will be used
func CommonLogHandler(logger *log.Logger, h http.Handler, newLogData func() CommonLogData) http.Handler {
	if logger == nil {
		logger = log.New(os.Stdout, "", 0)
	}
	if newLogData == nil {
		newLogData = NewDefaultCommonLogData
	}
	return &commonLogHandler{
		handler:    h,
		logger:     logger,
		newLogData: newLogData,
	}
}

// ServeHTTP logs the request and response data to Common Log Format
func (lh *commonLogHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	data := lh.newLogData()
	// decorate the writer so we can capture the status and size
	loggedWriter := response.LogResponseMetadata(w)
	// grab the data we need before passing it to ServeHTTP
	// since ServeHTTP could modify that data
	data.SetRequest(*req)
	lh.handler.ServeHTTP(loggedWriter, req)

	data.SetResponseWriter(*loggedWriter)
	lh.logger.Println(data.String())
}

type CommonLogData interface {
	SetRequest(http.Request)
	SetResponseWriter(response.ResponseMetadataLogger)
	String() string
}

type commonLogData struct {
	StartTime      time.Time
	RemoteAddr     string
	Username       string
	Method         string
	Uri            string
	Proto          string
	ResponseStatus int
	ResponseSize   int
}

func NewDefaultCommonLogData() CommonLogData {
	return &commonLogData{
		StartTime: time.Now(),
	}
}

func (this *commonLogData) SetRequest(req http.Request) {
	this.RemoteAddr = req.RemoteAddr
	this.Username = usernameFromReq(req)
	this.Method = req.Method
	this.Uri = req.RequestURI
	this.Proto = req.Proto
}

func (this *commonLogData) SetResponseWriter(rw response.ResponseMetadataLogger) {
	this.ResponseStatus = rw.Status
	this.ResponseSize = rw.Size
}

func (this *commonLogData) String() string {
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

// Type CommonLogUsingForwardedFor logs in CLF using
// the X-Forwarded-For header in place of the remote IP
type CommonLogUsingForwardedFor struct {
	CommonLogData
}

func NewCommonLogUsingForwardedFor() CommonLogData {
	return &CommonLogUsingForwardedFor{NewDefaultCommonLogData()}
}

func (this *CommonLogUsingForwardedFor) SetRequest(req http.Request) {
	req.RemoteAddr = req.Header.Get("X-Forwarded-For")
	this.CommonLogData.SetRequest(req)
}

// Type CommonLogUsingForwardedFor logs in CLF but
// excludes query params from the log
type CommonLogStrippingQueries struct {
	CommonLogData
}

func NewCommonLogStrippingQueries() CommonLogData {
	return &CommonLogStrippingQueries{NewDefaultCommonLogData()}
}

func (this *CommonLogStrippingQueries) SetRequest(req http.Request) {
	req.RequestURI = req.URL.Path
	this.CommonLogData.SetRequest(req)
}

// Extract username from the request url
func usernameFromReq(req http.Request) string {
	if req.URL.User != nil {
		if name := req.URL.User.Username(); name != "" {
			return name
		}
	}
	return "-"
}
