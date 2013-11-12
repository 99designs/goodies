package max_bytes

import (
	"net/http"
)

type maxBytesHandler struct {
	h        http.Handler
	maxBytes int64
}

func MaxBytesHandler(max int64, h http.Handler) http.Handler {
	return maxBytesHandler{h, max}
}

func (m maxBytesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, m.maxBytes)
	m.h.ServeHTTP(w, r)
}
