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
		r.URL.Path = "/changed"
		w.Header().Set("X-Path", r.URL.Path)
		w.Write([]byte("Testing 1 2 3"))
	})

	//	output := ""
	logw := bytes.NewBuffer(nil)
	// logw.Write([]byte("testing"))
	// log.Println(logw.String())
	testlog := log.New(logw, "", 0)

	ts := httptest.NewServer(CommonLogHandler(testlog, h))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/foo/bar")
	if err != nil {
		t.Fatal(err)
	}

	e := regexp.MustCompile(`^127.0.0.1:\d+ - - [.+] "GET /foo/bar HTTP/1.1" 200 13$`)
	g := logw.String()
	if e.MatchString(g) {
		t.Errorf("test 1: got %s, want %s", g, e)
	}
	res.Body.Close()
}
