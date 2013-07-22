// Package errorhandling provides An HTTP decorator which recovers from `panic`
package errorhandling

import (
	"log"
	"net/http"
)

type RecoveringHandler struct {
	http.Handler
	RecoveryHandler
}

type RecoveryHandler func(*http.Request, http.ResponseWriter, interface{})

func Decorate(delegate http.Handler, recoveryHandler RecoveryHandler) RecoveringHandler {
	if recoveryHandler == nil {
		recoveryHandler = DefaultRecoveryHandler
	}
	return RecoveringHandler{delegate, recoveryHandler}
}

func (lh RecoveringHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			lh.RecoveryHandler(r, rw, rec)
		}
	}()
	lh.Handler.ServeHTTP(rw, r)
}

func DefaultRecoveryHandler(r *http.Request, rw http.ResponseWriter, err interface{}) {
	log.Println("Fatal error: ", err)
	http.Error(rw, "Eeek", http.StatusInternalServerError)
}
