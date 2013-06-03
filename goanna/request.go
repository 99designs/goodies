package goanna

import (
	"io/ioutil"
	"net/http"
)

type Request struct {
	*http.Request
	session  Session
	bodyData []byte
	bodyRead bool
}

func (r *Request) CookieValue(name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return c.Value
}

func (r *Request) BodyData() []byte {
	if !r.bodyRead {
		if r.Body != nil {
			r.bodyData, _ = ioutil.ReadAll(r.Body)
		}
		r.bodyRead = true
	}

	return r.bodyData
}
