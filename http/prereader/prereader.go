// Package prereader avoids I/O timeouts in your controllers by getting the entire body ahead of time.
package prereader

import (
	"fmt"
	gioutil "github.com/99designs/goodies/ioutil"
	"io/ioutil"
	. "net/http"
)

type PrereadHandler struct {
	Handler
}

func Decorate(delegate Handler) Handler {
	return PrereadHandler{delegate}
}

func (t PrereadHandler) ServeHttp(w ResponseWriter, r *Request) {
	err := PrereadBody(r)
	if err == nil {
		t.Handler.ServeHTTP(w, r)
	} else {
		Error(w, "Error reading request: "+err.Error(), StatusRequestTimeout)
	}

}

func PrereadBody(r *Request) (result error) {
	defer func() {
		if rec := recover(); rec != nil {
			result = fmt.Errorf("Error prereading body: %+v\n", rec)
			if _, ok := rec.(error); ok {
				result = rec.(error)
			}
		}
	}()
	r.Body = gioutil.NewBufferedReadCloser(r.Body)
	_, result = ioutil.ReadAll(r.Body)
	return result
}
