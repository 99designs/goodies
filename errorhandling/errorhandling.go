// Package errorhandling provides An HTTP decorator which recovers from `panic`
package errorhandling

import (
	"bytes"
	requestLogging "github.com/99designs/goodies/http/log/request"
	responseLogging "github.com/99designs/goodies/http/log/response"
	"io"
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

type RecoveringHandler struct {
	http.Handler
	RecoveryHandler
	*log.Logger
}

type RecoveryHandler func(*http.Request, http.ResponseWriter, interface{})

func Decorate(delegate http.Handler, recoveryHandler RecoveryHandler) *RecoveringHandler {
	if recoveryHandler == nil {
		recoveryHandler = DefaultRecoveryHandler
	}
	return &RecoveringHandler{delegate, recoveryHandler, nil}
}

func (lh RecoveringHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var body requestLogging.LoggedRequestBody
	var writer responseLogging.LoggedResponseWriter

	if lh.Logger != nil {
		// This isn't free, so only do it if logging is enabled.
		body = requestLogging.LogRequestBody(r)
		writer = responseLogging.LogResponseWriter(rw)
	}

	defer func() {
		if rec := recover(); rec != nil {
			lh.RecoveryHandler(r, rw, rec)
			if lh.Logger != nil {
				serializedHeaders := bytes.Buffer{}
				r.Header.Write(&serializedHeaders)

				lh.Logger.Printf(LogFormat,
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

	lh.Handler.ServeHTTP(writer, r)
}

func DefaultRecoveryHandler(r *http.Request, rw http.ResponseWriter, err interface{}) {
	log.Println("Fatal error: ", err)
	log.Println(string(debug.Stack()))
	http.Error(rw, "Eeek", http.StatusInternalServerError)
}
