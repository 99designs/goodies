package goanna

import (
	"bytes"
	"log"
	"runtime/debug"
	"strings"
	"time"
)

const LogRequestTemplate = `
----------------------------------------------------------------------
%s

Url: %s
Method: %s
Timestamp: %s

Request Headers:
%s

Request Body:
%s

Stack trace:
%s
----------------------------------------------------------------------
`

var Logger *log.Logger

// LogRequest logs a request using the
func LogRequest(r *Request, v ...string) {
	serializedHeaders := bytes.Buffer{}
	r.Header.Write(&serializedHeaders)

	printf := log.Printf
	if Logger != nil {
		printf = Logger.Printf
	}

	printf(
		LogRequestTemplate,
		strings.Join(v, " "),
		r.URL.String(),
		r.Method,
		time.Now(),
		serializedHeaders.String(),
		string(r.BodyData()),
		debug.Stack(),
	)
}
