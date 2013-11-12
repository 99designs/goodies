package secure_headers

import (
	"net/http"
	_ "testing"
)

func ExampleDecorate() {
	var h http.Handler
	h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	h = Decorate(DefaultSettings, h) // Responses from h now include default security headers
}
