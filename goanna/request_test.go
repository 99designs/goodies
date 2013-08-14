package goanna

import (
	ghttp "github.com/99designs/goodies/http"
	"net/http"
	"testing"
)

type timeoutError struct{}

func (e *timeoutError) Error() string   { return "i/o timeout" }
func (e *timeoutError) Timeout() bool   { return true }
func (e *timeoutError) Temporary() bool { return true }

type timeoutReader struct{}

func (this timeoutReader) Read(p []byte) (n int, err error) {
	return 0, &timeoutError{}
}
func (this timeoutReader) Close() error {
	return nil
}

func TestIoTimeout(t *testing.T) {

	req, err := http.NewRequest("POST", "http://foo/bar", timeoutReader{})
	panicIf(err)
	r := Request{Request: req}

	defer func() {
		gotHttpErr := false
		if panicErr := recover(); panicErr != nil {
			_, gotHttpErr = panicErr.(ghttp.HttpError)
		}
		if !gotHttpErr {
			t.Error("Expected a HttpError")
		}
	}()

	r.BodyData()
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
