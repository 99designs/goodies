// Package panichandler provides An HTTP decorator which recovers from `panic`
package panichandler

import (
	"bytes"
	ghttp "github.com/99designs/goodies/http"
	responseLogging "github.com/99designs/goodies/http/log"
	gioutil "github.com/99designs/goodies/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

const LogFormat = `-----------------------------
*** Panic serving request ***
Url: %s
Method: %s
Timestamp: %s
****** Request Headers ******
%s
******* Request Body ********
%s
******* Response Body *******
%s
******* Panic details *******
%+v (%s)
******** Stack trace ********
%s
-----------------------------
`

type RecoveryHandlerFunc func(http.ResponseWriter, *http.Request, interface{})

type PanicHandler struct {
	handler  http.Handler
	recovery RecoveryHandlerFunc
	logger   *log.Logger
}

func Decorate(recovery RecoveryHandlerFunc, logger *log.Logger, delegate http.Handler) *PanicHandler {
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
	var requestBody gioutil.BufferedReadCloser
	var writer *responseLogging.LoggedResponseBodyWriter

	if lh.logger != nil {
		requestBody = gioutil.NewBufferedReadCloser(r.Body)
		r.Body = requestBody

		writer = responseLogging.LogResponseBody(rw)
	}

	defer func() {
		if panicErr := recover(); panicErr != nil {
			lh.recovery(rw, r, panicErr)

			if lh.logger != nil {
				serializedHeaders := bytes.Buffer{}
				r.Header.Write(&serializedHeaders)
				errDescription := ""
				if err, ok := panicErr.(error); ok {
					errDescription = err.Error()
				}

				lh.logger.Printf(LogFormat,
					r.URL.String(),
					r.Method,
					time.Now(),
					string(serializedHeaders.String()),
					string(requestBody.ReadAll()),
					string(writer.Output.String()),
					panicErr,
					errDescription,
					string(debug.Stack()),
				)
			}
		}
	}()

	lh.handler.ServeHTTP(writer, r)
}

func DefaultRecoveryHandler(rw http.ResponseWriter, r *http.Request, panicErr interface{}) {
	httperr, isHttpError := panicErr.(ghttp.HttpError)
	if isHttpError {
		httperror(rw, httperr.StatusCode())
	} else {
		httperror(rw, http.StatusInternalServerError)
	}
}

func httperror(rw http.ResponseWriter, statuscode int) {
	http.Error(rw, http.StatusText(statuscode), statuscode)
}
