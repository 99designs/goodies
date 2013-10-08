package goanna

import (
	"net/http"
	"testing"
)

func TestEncode(t *testing.T) {
	cs := NewCookieSigner([]byte("TESTSECRET"))

	v := "foo.bar"
	v1 := cs.EncodeValue(v)
	v2, err2 := cs.DecodeValue(v1)
	_, err3 := cs.DecodeValue(v1 + "blah")

	if v == v1 {
		t.Error("Values should not be equal")
	}
	if v != v2 {
		t.Errorf("Expected %s, got %s", v, v2)
		t.Error(err2)
	}
	if err3 == nil {
		t.Error("Should have an error")
	}

	val := "Rm9av.shZ0lvgf-z5Rh7sfuzRZatPN0Y1LCJ54GCf0Tq8X0f8="
	c := http.Cookie{Name: "userid", Value: val}
	cs.EncodeCookie(&c)

	if val == c.Value {
		t.Error("Values should not be equal")
	}

	cs.DecodeCookie(&c)
	if v != v2 {
		t.Error("Values should be equal")
	}
}
