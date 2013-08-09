package log

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func setupTestServer(format string) (*httptest.Server, *bytes.Buffer) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Method = "CHANGED"
		w.Header().Set("X-Path", r.URL.Path)
		w.Write([]byte("Testing 1 2 3"))
	})

	logw := bytes.NewBuffer(nil)
	testlog := log.New(logw, "", 0)

	ts := httptest.NewServer(CommonLogHandler(testlog, format, h))
	return ts, logw
}

func TestCommonLogHandler(t *testing.T) {
	ts, logw := setupTestServer(DefaultFormat)
	defer ts.Close()
	req, err := http.NewRequest("GET", ts.URL+"/foo/bar", nil)
	panicIf(err)
	res, err := http.DefaultClient.Do(req)
	panicIf(err)

	e := regexp.MustCompile(`^127.0.0.1:\d+ - - .+ "GET /foo/bar HTTP/1.1" 200 13`)
	g := logw.String()
	if !e.MatchString(g) {
		t.Errorf("test 1: got '%s', want '%s'", g, e)
	}
	res.Body.Close()
}

func TestXForwardedLogger(t *testing.T) {
	ts, logw := setupTestServer(DefaultFormatUsingXForwardedFor)
	req, err := http.NewRequest("GET", ts.URL+"/foo/bar", nil)
	panicIf(err)
	req.Header.Add("X-Forwarded-For", "xff_value")
	res, err := http.DefaultClient.Do(req)
	panicIf(err)
	defer res.Body.Close()

	e := regexp.MustCompile(`^xff_value `)
	g := logw.String()
	if !e.MatchString(g) {
		t.Errorf("test 1: got '%s', want '%s'", g, e)
	}
}

func TestStripQuery(t *testing.T) {
	ts, logw := setupTestServer(DefaultFormatStrippingQueries)
	req, err := http.NewRequest("GET", ts.URL+"/foo/bar?baz=qux", nil)
	panicIf(err)
	res, err := http.DefaultClient.Do(req)
	panicIf(err)
	defer res.Body.Close()

	e := regexp.MustCompile(`"GET /foo/bar HTTP/1.1"`)
	g := logw.String()
	if !e.MatchString(g) {
		t.Errorf("test 1: got '%s', want '%s'", g, e)
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
