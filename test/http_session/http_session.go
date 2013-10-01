package http_session

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

func GetRequest(geturl string, cookies []http.Cookie) *http.Request {
	return request("GET", geturl, cookies, nil)
}

func PostRequest(posturl string, cookies []http.Cookie, postdata url.Values) *http.Request {
	return request("POST", posturl, cookies, postdata)
}

func request(method, url string, cookies []http.Cookie, data url.Values) *http.Request {
	var body io.Reader
	if data != nil {
		body = strings.NewReader(data.Encode())
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	for _, c := range cookies {
		req.AddCookie(&c)
	}
	if data != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return req
}
