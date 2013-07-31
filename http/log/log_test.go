package log

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestCommonLogHandler(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Method = "CHANGED"
		w.Header().Set("X-Path", r.URL.Path)
		w.Write([]byte("Testing 1 2 3"))
	})

	logw := bytes.NewBuffer(nil)
	testlog := log.New(logw, "", 0)

	ts := httptest.NewServer(CommonLogHandler(testlog, h, nil))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/foo/bar")
	if err != nil {
		t.Fatal(err)
	}

	e := regexp.MustCompile(`^127.0.0.1:\d+ - - .+ "GET /foo/bar HTTP/1.1" 200 13`)
	g := logw.String()
	if !e.MatchString(g) {
		t.Errorf("test 1: got '%s', want '%s'", g, e)
	}
	res.Body.Close()
}

func TestXForwardedLogger(t *testing.T) {
	l := NewCommonLogUsingForwardedFor()
	r, err := http.NewRequest("GET", "/foo", nil)
	r.RequestURI = "/foo"
	if err != nil {
		t.Error(err)
		return
	}
	r.Header.Set("X-Forwarded-For", "10.0.0.1")
	l.SetRequest(*r)

	e := regexp.MustCompile(`^10.0.0.1`)
	g := l.String()
	if !e.MatchString(g) {
		t.Errorf("test 1: got '%s', want '%s'", g, e)
	}
}

func TestStripQuery(t *testing.T) {
	l := NewCommonLogStrippingQueries()
	r, err := http.NewRequest("GET", "/foo?bar=baz", nil)
	r.RequestURI = "/foo?bar=baz"
	if err != nil {
		t.Error(err)
		return
	}
	l.SetRequest(*r)

	e := regexp.MustCompile(`"GET /foo HTTP/1.1"`)
	g := l.String()
	if !e.MatchString(g) {
		t.Errorf("test 1: got '%s', want '%s'", g, e)
	}
}
