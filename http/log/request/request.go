package request

import (
	"bytes"
	"io"
	"net/http"
)

type LoggedRequestBody struct {
	orig   io.ReadCloser
	Output bytes.Buffer
}

func LogRequestBody(r *http.Request) LoggedRequestBody {
	output := bytes.Buffer{}
	result := LoggedRequestBody{
		orig:   r.Body,
		Output: output,
	}
	r.Body = result
	return result
}

func (this LoggedRequestBody) Read(p []byte) (n int, err error) {
	n, err = this.orig.Read(p)
	_, _ = this.Output.Write(p[:n])
	return n, err
}

func (this LoggedRequestBody) Close() error {
	return this.orig.Close()
}
