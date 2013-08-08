// Package panichandler provides An HTTP decorator which recovers from `panic`
package panichandler

import (
	"bytes"
	requestLogging "github.com/99designs/goodies/http/log/request"
	responseLogging "github.com/99designs/goodies/http/log/response"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

const LogFormat = `---------------------------------------
****** Panic serving request ******
Url: %s
Method: %s
Timestamp: %s
****** Headers ******
%s
****** Request Body ******
%s
****** Response Body ******
%s
****** Panic details ******
%+v
****** Stack trace ******
%s
---------------------------------------
`

type PanicHandler struct {
	handler  http.Handler
	recovery http.HandlerFunc
	logger   *log.Logger
}

func Decorate(delegate http.Handler, recovery http.HandlerFunc, logger *log.Logger) *PanicHandler {
	if recovery == nil {
		recovery = DefaultRecoveryHandler
	}
	return &PanicHandler{
		handler:  delegate,
		recovery: recovery,
		logger:   logger,
	}
}

func (lh PanicHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var body requestLogging.LoggedRequestBody
	var writer *responseLogging.LoggedResponseBodyWriter

	if lh.logger != nil {
		// This isn't free, so only do it if logging is enabled.
		body = requestLogging.LogRequestBody(r)
		writer = responseLogging.LogResponseBody(rw)
	}

	defer func() {
		if rec := recover(); rec != nil {
			lh.recovery(rw, r)

			if lh.logger != nil {
				serializedHeaders := bytes.Buffer{}
				r.Header.Write(&serializedHeaders)

				lh.logger.Printf(LogFormat,
					r.URL.String(),
					r.Method,
					time.Now(),
					string(serializedHeaders.String()),
					string(body.Output.String()),
					string(writer.Output.String()),
					rec,
					string(debug.Stack()),
				)
			}
		}
	}()

	lh.handler.ServeHTTP(writer, r)
}

func DefaultRecoveryHandler(rw http.ResponseWriter, r *http.Request) {
	http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
}
