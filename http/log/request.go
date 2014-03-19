package log

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("** DEBUG: Incoming Request ---------------")
		fmt.Printf("%s %s\n", r.Method, r.URL.String())
		r.Header.Write(os.Stdout)

		// read body but leave r.Body available for reading
		body, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewReader(body))

		fmt.Println("\n", string(body))
		fmt.Println("** ---------------------------------------")
		h.ServeHTTP(HttpLogger{rw}, r)
	})
}

type HttpLogger struct {
	rw http.ResponseWriter
}

func (hl HttpLogger) Header() http.Header {
	return hl.rw.Header()
}

func (hl HttpLogger) Write(content []byte) (int, error) {
	fmt.Println("** DEBUG: Writing response: --------------")
	fmt.Println(string(content))
	fmt.Println("** ---------------------------------------")
	return hl.rw.Write(content)
}

func (hl HttpLogger) WriteHeader(code int) {
	fmt.Println("** DEBUG: Writing header: ", code)
	hl.rw.WriteHeader(code)
}
