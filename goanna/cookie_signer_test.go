package goanna

import (
	e "github.com/99designs/goodies/test/expect"
	"net/http"
	"testing"
)

func TestEncode(t *testing.T) {
	c := NewCookieSigner("TESTSECRET")

	e.Expect(t, "Cookie value", c.encodeValue("Foo"), "Rm9v.shZ0lvgf-z5Rh7sfuzRZatPN0Y1LCJ54GCf0Tq8X0f8=")
	e.Expect(t, "Full cookie string", c.EncodeCookie(http.Cookie{Name: "userid", Value: "Foo"}).String(), "userid=Rm9v.shZ0lvgf-z5Rh7sfuzRZatPN0Y1LCJ54GCf0Tq8X0f8=")

	cookie, _ := c.DecodeCookie(http.Cookie{Name: "userid", Value: "Rm9v.shZ0lvgf-z5Rh7sfuzRZatPN0Y1LCJ54GCf0Tq8X0f8="})
	e.Expect(t, "Decode valid cookie string", cookie.Value, "Foo")

	_, err := c.DecodeCookie(http.Cookie{Name: "userid", Value: "Rm9av.shZ0lvgf-z5Rh7sfuzRZatPN0Y1LCJ54GCf0Tq8X0f8="})
	e.Expect(t, "Decode invalid cookie string", err.Error(), "illegal base64 data at input byte 4")
}
