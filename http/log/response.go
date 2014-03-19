package log

import (
	"bytes"
	"io"
	"net/http"
)

type LoggedResponseBodyWriter struct {
	rw     http.ResponseWriter
	Output bytes.Buffer
	multi  io.Writer
}

func LogResponseBody(rw http.ResponseWriter) *LoggedResponseBodyWriter {
	output := bytes.Buffer{}
	return &LoggedResponseBodyWriter{
		rw:     rw,
		Output: output,
		multi:  io.MultiWriter(&output, rw),
	}
}

func (this *LoggedResponseBodyWriter) Header() http.Header {
	return this.rw.Header()
}

func (this *LoggedResponseBodyWriter) Write(p []byte) (int, error) {
	return this.multi.Write(p)
}

func (this *LoggedResponseBodyWriter) WriteHeader(h int) {
	this.rw.WriteHeader(h)
}

type ResponseWriterWatcher struct {
	rw     http.ResponseWriter
	Status int
	Size   int
}

func WatchResponseWriter(rw http.ResponseWriter) *ResponseWriterWatcher {
	return &ResponseWriterWatcher{
		rw: rw,
	}
}

func (l *ResponseWriterWatcher) Header() http.Header {
	return l.rw.Header()
}

func (l *ResponseWriterWatcher) Write(b []byte) (int, error) {
	if l.Status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		l.Status = http.StatusOK
	}
	size, err := l.rw.Write(b)
	l.Size += size
	return size, err
}

func (l *ResponseWriterWatcher) WriteHeader(s int) {
	l.rw.WriteHeader(s)
	l.Status = s
}
