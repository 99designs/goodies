package goanna

import (
	"encoding/json"
	errorLogging "github.com/99designs/goodies/http/log/error"
	"log"
	"net/http"
	"time"
)

const (
	HEADER_CACHE_CONTROL = "Cache-Control"
	HEADER_CACHE_NOCACHE = "no-cache"
)

type Response interface {
	Send(http.ResponseWriter)
	SetCookie(http.Cookie)
	ClearCookie(string)
	GetStatusCode() int
	Headers() http.Header
}

type OkResponse struct {
	Header  http.Header
	Cookies []http.Cookie
	Content []byte
}

func NewResponse() *OkResponse {
	return &OkResponse{
		Header:  http.Header{},
		Cookies: make([]http.Cookie, 0),
		Content: make([]byte, 0),
	}
}

func NewStringResponse(content string) Response {
	response := NewResponse()
	response.Content = []byte(content)
	return response
}

func (r *OkResponse) Headers() http.Header {
	return r.Header
}

func (r *OkResponse) SetCookie(cookie http.Cookie) {
	r.Cookies = append(r.Cookies, cookie)
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
	r.Content = content
}

func (r *OkResponse) SendHeaders(w http.ResponseWriter) {
	for _, cookie := range r.Cookies {
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
	w.Write([]byte(r.Content))
}

func (r *OkResponse) SetNoCache() {
	r.Header.Set(HEADER_CACHE_CONTROL, HEADER_CACHE_NOCACHE)
}

func (r *OkResponse) GetStatusCode() int {
	return 200
}

type ErrorResponse struct {
	OkResponse
	Message string
	Code    int
}

func NewErrorResponse(message string, code int) *ErrorResponse {
	// Log a minimal error
	errorLogging.LogError(
		nil,
		nil,
		message,
		"",
		"",
	)

	response := ErrorResponse{
		OkResponse: *NewResponse(),
		Message:    message,
		Code:       code,
	}

	return &response
}

func (r *ErrorResponse) Send(w http.ResponseWriter) {
	r.SendHeaders(w)
	http.Error(w, r.Message, r.Code)
	r.SetNoCache()
}

func (r *ErrorResponse) GetStatusCode() int {
	return r.Code
}

type RedirectResponse struct {
	OkResponse
	urlStr string
	code   int
}

func NewRedirectResponse(urlStr string, code int) *RedirectResponse {
	response := RedirectResponse{
		OkResponse: *NewResponse(),
		urlStr:     urlStr,
		code:       code,
	}

	return &response
}

func (r *RedirectResponse) Send(w http.ResponseWriter) {
	w.Header().Set("Location", r.urlStr)
	r.SendHeaders(w)
	w.WriteHeader(r.code)
}

func (r *RedirectResponse) Target() string {
	return r.urlStr
}
func (r *RedirectResponse) Code() int {
	return r.code
}

func NewJsonResponse(data interface{}) *OkResponse {
	json, err := json.Marshal(data)
	if err != nil {
		log.Panicln("Bad data for json marshalling")
	}

	response := NewResponse()
	response.Content = []byte(json)
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
	response.Content = body
	response.Header.Set("Content-Type", "application/javascript")

	return response
}
