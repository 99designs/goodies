package goanna

import (
	e "github.com/99designs/goodies/test/expect"
	"testing"
)

func TestRedirectGetters(t *testing.T) {
	r := RedirectResponse{urlStr: "Foo", code: 1234}

	e.Expect(t, "Target Getter", r.Target(), "Foo")
	e.Expect(t, "Code Getter", r.Code(), 1234)
}
