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
	var err error
	if !r.bodyRead {
		if r.Body != nil {
			r.bodyData, err = ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
		}
		r.bodyRead = true
	}

	return r.bodyData
}

func (r *Request) QueryValue(key string) string {
	return r.URL.Query().Get(key)
}
