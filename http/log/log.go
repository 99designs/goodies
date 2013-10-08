// Package log decorates any http handler with a
// Common Log Format logger
package log

import (
	"bytes"
	"github.com/99designs/goodies/http/log/response"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
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

	t := DefaultFormatter
	if templ != "" {
		t = template.Must(template.New("").Parse(templ))
	}

	return &commonLogHandler{
		handler:  h,
		logger:   logger,
		template: t,
	}
}

// ServeHTTP logs the request and response data to Common Log Format
func (lh *commonLogHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// decorate the writer so we can capture the status and size
	loggedWriter := response.LogResponseMetadata(w)
	// copy the request struct into loggableData so it cannot be modified
	loggableData := LoggableData{time.Now(), *req, loggedWriter}
	// Fall through to the actual http handler
	lh.handler.ServeHTTP(loggedWriter, req)

	logOutput := bytes.Buffer{}
	err := lh.template.Execute(&logOutput, loggableData)
	if err != nil {
		panic(err)
	}
	lh.logger.Println(logOutput.String())
}

const CommonLogFormat = `{{ or .Req.RemoteAddr '-' }} {{ or .UrlUsername '-' }} - [{{ or .CLFFormatStartTime '-' }}] "{{ or .Req.Method '-' }} {{ or .Req.RequestURI '-' }} {{ or .Req.Proto '-' }}" {{ or .Rw.Status '-' }} {{ or .Rw.Size '-' }}`
const CommonLogFormatUsingXForwardedFor = `{{ or ( .Req.Header.Get "X-Forwarded-For" ) '-' }} {{ or .UrlUsername '-' }} - [{{ or .CLFFormatStartTime '-' }}] "{{ or .Req.Method '-' }} {{ or .Req.RequestURI '-' }} {{ or .Req.Proto '-' }}" {{ or .Rw.Status '-' }} {{ or .Rw.Size '-' }}`
const CommonLogFormatStrippingQueries = `{{ or .Req.RemoteAddr '-' }} {{ or .UrlUsername '-' }} - [{{ or .CLFFormatStartTime '-' }}] "{{ or .Req.Method '-' }} {{ or .Req.URL.Path '-' }} {{ or .Req.Proto '-' }}" {{ or .Rw.Status '-' }} {{ or .Rw.Size '-' }}`

var DefaultFormatter *template.Template

func init() {
	DefaultFormatter = template.Must(template.New("").Parse(CommonLogFormat))
}

type LoggableData struct {
	StartTime time.Time
	Req       http.Request
	Rw        *response.ResponseMetadataLogger
}

// Format the start time
func (this LoggableData) CLFFormatStartTime() string {
	return this.StartTime.Format("02/Jan/2006:15:04:05 -0700")
}

// Extract username from the request url
func (this LoggableData) UrlUsername() string {
	if this.Req.URL.User != nil {
		if name := this.Req.URL.User.Username(); name != "" {
			return name
		}
	}
	return "-"
}
