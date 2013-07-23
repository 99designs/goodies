package logged_response_writer

import (
	"bytes"
	"io"
	"net/http"
)

type LoggedResponseWriter struct {
	rw     http.ResponseWriter
	Output bytes.Buffer
	multi  io.Writer
}

func LogResponseWriter(rw http.ResponseWriter) LoggedResponseWriter {
	output := bytes.Buffer{}
	return LoggedResponseWriter{
		rw:     rw,
		Output: output,
		multi:  io.MultiWriter(&output, rw),
	}
}

func (this LoggedResponseWriter) Header() http.Header {
	rw := this.rw
	return rw.Header()
}

func (this LoggedResponseWriter) Write(p []byte) (int, error) {
	return this.multi.Write(p)
}

func (this LoggedResponseWriter) WriteHeader(h int) {
	this.rw.WriteHeader(h)
}
