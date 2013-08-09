package panichandler

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

const PANICMSG = "panicmsg"
const BODYMSG = "panicmsg"

type pannicker struct{}

func setupTestServer() (*httptest.Server, *bytes.Buffer) {

	logw := bytes.NewBuffer(nil)
	testlog := log.New(logw, "", 0)

	panicker := Decorate(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(PANICMSG)
		}), nil, testlog)

	return httptest.NewServer(panicker), logw
}

func TestPanicHandler(t *testing.T) {
	path := "/foo/bar"
	ts, logw := setupTestServer()
	defer ts.Close()
	req, err := http.NewRequest("POST", ts.URL+path, bytes.NewReader([]byte(BODYMSG)))
	panicIf(err)
	res, err := http.DefaultClient.Do(req)
	panicIf(err)

	fmt.Println(res.Body)
	g := logw.String()
	fmt.Println(g)

	expect1 := regexp.MustCompile(PANICMSG)
	if !expect1.MatchString(g) {
		t.Errorf("test 1: got '%s', want '%s'", g, expect1)
	}

	expect2 := regexp.MustCompile(BODYMSG)
	if !expect2.MatchString(g) {
		t.Errorf("test 2: got '%s', want '%s'", g, expect2)
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
