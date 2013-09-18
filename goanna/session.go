package goanna

import (
	"time"
)

type SessionHandler interface {
	GetSession(*Request) Session
}

type Session interface {
	GetId() string
	Get(string) string
	Set(string, string)
	SetMaxAge(time.Duration)
	HasExpired() bool
	Clear()
	WriteToResponse(Response)
}
