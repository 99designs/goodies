package goanna

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const (
	HEADER_CACHE_CONTROL = "Cache-Control"
	HEADER_CACHE_NOCACHE = "no-cache"
)

// Response is a http response that can be built up
type Response interface {
	Send(http.ResponseWriter)
	SetCookie(http.Cookie)
	ClearCookie(string)
	Cookies() []http.Cookie
	Content() []byte
	StatusCode() int
	SetStatusCode(int)
	Headers() http.Header
}

type OkResponse struct {
	Header  http.Header
	cookies []http.Cookie
	content []byte
	Code    int
}

func NewResponse() *OkResponse {
	return &OkResponse{
		Header:  http.Header{},
		cookies: make([]http.Cookie, 0),
		content: make([]byte, 0),
		Code:    http.StatusOK,
	}
}

func NewStringResponse(content string) Response {
	response := NewResponse()
	response.content = []byte(content)
	return response
}

func (r *OkResponse) Cookies() []http.Cookie {
	return r.cookies
}

func (r *OkResponse) Headers() http.Header {
	return r.Header
}

func (r *OkResponse) Content() []byte {
	return r.content
}

func (r *OkResponse) SetCookie(cookie http.Cookie) {
	r.cookies = append(r.cookies, cookie)
}

func (r *OkResponse) ClearCookie(name string) {
	c := http.Cookie{}
	c.Name = name
	c.Value = "none"
	c.MaxAge = -1
	c.Expires = time.Unix(1, 0)
	r.SetCookie(c)
}

func (r *OkResponse) SetContent(content []byte) {
	r.content = content
}

func (r *OkResponse) SendHeaders(w http.ResponseWriter) {
	for _, cookie := range r.cookies {
		http.SetCookie(w, &cookie)
	}

	keys := make([]string, 0, len(r.Header))
	for k := range r.Header {
		keys = append(keys, k)
	}
	for _, k := range keys {
		for _, v := range r.Header[k] {
			w.Header().Add(k, v)
		}
	}
}

func (r *OkResponse) Send(w http.ResponseWriter) {
	r.SendHeaders(w)
	w.WriteHeader(r.Code)
	w.Write([]byte(r.content))
}

func (r *OkResponse) SetNoCache() {
	r.Header.Set(HEADER_CACHE_CONTROL, HEADER_CACHE_NOCACHE)
}

func (r *OkResponse) StatusCode() int {
	return r.Code
}

func (r *OkResponse) SetStatusCode(code int) {
	r.Code = code
}

type ErrorResponse struct {
	*OkResponse
	Message string
}

func NewErrorResponse(message string, code int) *ErrorResponse {
	response := ErrorResponse{
		OkResponse: NewResponse(),
		Message:    message,
	}
	response.SetStatusCode(code)

	return &response
}

func (r *ErrorResponse) Send(w http.ResponseWriter) {
	r.SendHeaders(w)
	http.Error(w, r.Message, r.Code)
	r.SetNoCache()
}

type RedirectResponse struct {
	*OkResponse
	urlStr string
}

func NewPermanentRedirectResponse(urlStr string) *RedirectResponse {
	response := RedirectResponse{
		OkResponse: NewResponse(),
		urlStr:     urlStr,
	}
	response.SetStatusCode(http.StatusMovedPermanently)

	return &response
}

func NewRedirectResponse(urlStr string) *RedirectResponse {
	response := RedirectResponse{
		OkResponse: NewResponse(),
		urlStr:     urlStr,
	}
	response.SetStatusCode(http.StatusFound)

	return &response
}

func (r *RedirectResponse) Send(w http.ResponseWriter) {
	w.Header().Set("Location", r.urlStr)
	r.SendHeaders(w)
	w.WriteHeader(r.Code)
}

func (r *RedirectResponse) Target() string {
	return r.urlStr
}

func NewJsonResponse(data interface{}) *OkResponse {
	json, err := json.Marshal(data)
	if err != nil {
		log.Panicln("Bad data for json marshalling")
	}

	response := NewResponse()
	response.content = []byte(json)
	response.Header.Set("Content-Type", "application/json")

	return response
}

func NewJsonpResponse(callback string, data interface{}) *OkResponse {
	json, err := json.Marshal(data)
	if err != nil {
		log.Panicln("Bad data for json marshalling")
	}
	body := append([]byte(callback+"("), json...)
	body = append(body, []byte(");")...)

	response := NewResponse()
	response.content = body
	response.Header.Set("Content-Type", "application/javascript")

	return response
}
