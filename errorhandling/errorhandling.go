// Package errorhandling provides An HTTP decorator which recovers from `panic`
package errorhandling

import (
	"bytes"
	"github.com/99designs/goodies/http/logged_request_body"
	"github.com/99designs/goodies/http/logged_response_writer"
	"io"
	"log"
	"net/http"
	"runtime/debug"
)

const LogFormat = `---------------------------------------
****** Panic serving request ******
Url: %s
Method: %s
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

func (lh *RecoveringHandler) LogTo(out io.Writer) {
	lh.Logger = log.New(out, "", 0)
}

func (lh RecoveringHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var body logged_request_body.LoggedRequestBody
	var writer logged_response_writer.LoggedResponseWriter

	if lh.Logger != nil {
		// This isn't free, so only do it if logging is enabled.
		body = logged_request_body.LogRequestBody(r)
		writer = logged_response_writer.LogResponseWriter(rw)
	}

	defer func() {
		if rec := recover(); rec != nil {
			lh.RecoveryHandler(r, rw, rec)
			log.Println(lh)
			log.Println(string(debug.Stack()))
			if lh.Logger != nil {
				serializedHeaders := bytes.Buffer{}
				r.Header.Write(&serializedHeaders)

				lh.Logger.Printf(LogFormat,
					r.URL.String(),
					r.Method,
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
	http.Error(rw, "Eeek", http.StatusInternalServerError)
}
