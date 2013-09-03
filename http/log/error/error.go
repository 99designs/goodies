// Package error provides logging for http errors
package error

import (
	"bytes"
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

func LogError(l *log.Logger, r *http.Request, err interface{}, reqBody string, respBody string) {
	// Setup logger
	logFunc := l.Printf
	if l == nil {
		logFunc = log.Printf
	}

	// Get error string if possible
	errDescription := ""
	if err, ok := err.(error); ok {
		errDescription = err.Error()
	}

	// Read url/method/headers if request available
	url := ""
	method := ""
	serializedHeaders := bytes.Buffer{}
	if r != nil {
		url = r.URL.String()
		method = r.Method
		r.Header.Write(&serializedHeaders)
	}

	// Log the details we've captured
	logFunc(LogFormat,
		url,
		method,
		time.Now(),
		string(serializedHeaders.String()),
		reqBody,
		respBody,
		err,
		errDescription,
		string(debug.Stack()),
	)
}
