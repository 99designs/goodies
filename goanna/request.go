package goanna

import (
	ghttp "github.com/99designs/goodies/http"
	"io/ioutil"
	"net"
	"net/http"
)

// Request decorates a http.Request to add helper methods
type Request struct {
	*http.Request
	session  Session
	bodyData []byte
	bodyRead bool
}

// CookieValue returns the value of the named cookie
func (r *Request) CookieValue(name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return c.Value
}

// BodyData returns the full request body
func (r *Request) BodyData() []byte {
	var err error
	if !r.bodyRead {
		if r.Body != nil {
			r.bodyData, err = ioutil.ReadAll(r.Body)
			if err != nil {
				// catch i/o timeout errors
				neterr, isNetError := err.(net.Error)
				if isNetError && neterr.Timeout() {
					panic(ghttp.NewHttpError(err, http.StatusRequestTimeout))
				} else {
					panic(err)
				}
			}
		}
		r.bodyRead = true
	}

	return r.bodyData
}

// QueryValue returns the value in the GET query string
func (r *Request) QueryValue(key string) string {
	return r.URL.Query().Get(key)
}

func (r *Request) QueryValueOrDefault(key string, def string) string {
	val := r.URL.Query().Get(key)
	if val == "" {
		val = def
	}

	return val
}
