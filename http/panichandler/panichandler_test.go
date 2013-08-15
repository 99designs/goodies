package panichandler

import (
	"bytes"
	"errors"
	ghttp "github.com/99designs/goodies/http"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

const PANICMSG = "panicmsg"
const BODYMSG = "panicmsg"

type pannicker struct{}

func TestPanicHandler(t *testing.T) {
	logw := bytes.NewBuffer(nil)
	testlog := log.New(logw, "", 0)
	delegate := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(PANICMSG)
	})
	panicker := Decorate(
		nil, testlog, delegate)

	ts := httptest.NewServer(panicker)
	defer ts.Close()

	path := "/foo/bar"
	req, err := http.NewRequest("POST", ts.URL+path, bytes.NewReader([]byte(BODYMSG)))
	panicIf(err)
	_, err = http.DefaultClient.Do(req)
	panicIf(err)
	g := logw.String()

	expect1 := regexp.MustCompile(PANICMSG)
	if !expect1.MatchString(g) {
		t.Errorf("test 1: got '%s', want '%s'", g, expect1)
	}

	expect2 := regexp.MustCompile(BODYMSG)
	if !expect2.MatchString(g) {
		t.Errorf("test 2: got '%s', want '%s'", g, expect2)
	}
}

func TestHttpErrorHandler(t *testing.T) {
	delegate := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(ghttp.HttpError{errors.New("i/o timeout"), http.StatusRequestTimeout})
	})
	h := Decorate(nil, nil, delegate)

	ts := httptest.NewServer(h)
	defer ts.Close()

	req, err := http.NewRequest("POST", ts.URL, bytes.NewReader([]byte("blah")))
	panicIf(err)
	res, err := http.DefaultClient.Do(req)
	panicIf(err)

	if res.StatusCode != http.StatusRequestTimeout {
		t.Errorf("Expected %d, got %d", http.StatusRequestTimeout, res.StatusCode)
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
