package ratelimiter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDecorator(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Path", r.URL.Path)
		w.Write([]byte("Testing 1 2 3"))
	})

	rl := func(a *http.Request) (ok bool, retryAfter time.Duration) {
		return false, time.Minute
	}

	testlog := Decorate(h, rl)

	ts := httptest.NewServer(testlog)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/foo/bar")
	if err != nil {
		t.Fatal(err)
	}
	if retryAfter := res.Header.Get("Retry-After"); retryAfter != "60" {
		t.Fatalf("Expected retry-after to be '60' but it was '%s'", retryAfter)
	}
	res.Body.Close()
}
