package goanna_test

import (
	"fmt"
	"github.com/99designs/goodies/goanna"
	"net/http"
)

func ExampleCookieSigner() {
	signer := goanna.NewCookieSigner([]byte("secret"))

	cookie := http.Cookie{
		Name:  "foo",
		Value: "bar",
	}

	fmt.Println(cookie.Value) // bar

	signer.EncodeCookie(&cookie)
	fmt.Println(cookie.Value) // aMcN2wzx4XLp9w3CPrwNb6PtTzECzkMPIiEfDqVDk4k=.bar

	signer.DecodeCookie(&cookie)
	fmt.Println(cookie.Value) // bar
}
