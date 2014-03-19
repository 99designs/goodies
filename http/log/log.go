// Package log is for http logging
package log

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

const (
	// templates for logging
	CommonLogFormat                   = `{{ or .Request.RemoteAddr "-" }} {{ or .UrlUsername "-" }} - [{{ or (.ClfTime .Start) "-" }}] "{{ or .Request.Method "-" }} {{ or .Request.RequestURI "-" }} {{ or .Request.Proto "-" }}" {{ or .Response.Status "-" }} {{ or .Response.Size "-" }}`
	CommonLogFormatUsingXForwardedFor = `{{ or ( .Request.Header.Get "X-Forwarded-For" ) "-" }} {{ or .UrlUsername "-" }} - [{{ or (.ClfTime .Start) "-" }}] "{{ or .Request.Method "-" }} {{ or .Request.RequestURI "-" }} {{ or .Request.Proto "-" }}" {{ or .Response.Status "-" }} {{ or .Response.Size "-" }}`
	CommonLogFormatStrippingQueries   = `{{ or .Request.RemoteAddr "-" }} {{ or .UrlUsername "-" }} - [{{ or (.ClfTime .Start) "-" }}] "{{ or .Request.Method "-" }} {{ or .Request.URL.Path "-" }} {{ or .Request.Proto "-" }}" {{ or .Response.Status "-" }} {{ or .Response.Size "-" }}`
)

type commonLogHandler struct {
	handler  http.Handler
	logger   *log.Logger
	template *template.Template
}

// CommonLogHandler returns a handler that serves HTTP requests
// If a logger is not provided, stdout will be used
func CommonLogHandler(logger *log.Logger, templ string, h http.Handler) http.Handler {
	if logger == nil {
		logger = log.New(os.Stdout, "", 0)
	}
	if templ == "" {
		templ = CommonLogFormat
	}

	return &commonLogHandler{
		handler:  h,
		logger:   logger,
		template: template.Must(template.New("").Parse(templ)),
	}
}

// ServeHTTP logs the request and response data to Common Log Format
func (lh *commonLogHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// take a copy of the request (so it's immutable)
	// watch the writer and capture the status and size
	logData := LogData{
		Request:  *req,
		Response: WatchResponseWriter(w),
	}

	logData.Start = time.Now()
	lh.handler.ServeHTTP(logData.Response, req)
	logData.End = time.Now()

	logOutput := bytes.Buffer{}
	err := lh.template.Execute(&logOutput, logData)
	if err != nil {
		panic(err)
	}

	lh.logger.Println(logOutput.String())
}

type LogData struct {
	Start    time.Time
	End      time.Time
	Request  http.Request
	Response *ResponseWriterWatcher
}

// Format time in CLF format
func (this LogData) ClfTime(t time.Time) string {
	return t.Format("02/Jan/2006:15:04:05 -0700")
}

// Extract username from the request url
func (this LogData) UrlUsername() string {
	if this.Request.URL.User != nil {
		return this.Request.URL.User.Username()
	}
	return ""
}
