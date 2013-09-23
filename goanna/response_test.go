package goanna

import (
	"testing"
)

func TestRedirectGetters(t *testing.T) {
	r := NewRedirectResponse("Foo")

	if r.Target() != "Foo" {
		t.Error("Target Getter")
	}
	if r.StatusCode() != 302 {
		t.Error("Code Getter")
	}
}
