package session

import (
	"github.com/99designs/goodies/goanna"
	"net/http"
	"testing"
	"time"
)

func TestSessionRoundTrip(t *testing.T) {
	finder := NewSignedCookieSessionFinder("cookiesession", []byte("secret"), 60*time.Minute, false)

	request1, _ := http.NewRequest("GET", "http://example.org/", nil)
	session1 := finder(&goanna.Request{Request: request1})
	session1.Set("foo", "bar")
	response1 := goanna.NewResponse()
	session1.WriteToResponse(response1)

	request2, _ := http.NewRequest("GET", "http://example.org/", nil)
	for _, c := range response1.Cookies() {
		newc := c
		request2.AddCookie(&newc)
	}
	session2 := finder(&goanna.Request{Request: request2})

	expected := "bar"
	actual := session2.Get("foo")
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
