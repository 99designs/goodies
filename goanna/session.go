package goanna

import (
	"time"
)

type SessionStore interface {
	GetSession(*Request) Session
}

type Session interface {
	Get(string) string
	Set(string, string)
	Expiry() time.Time
	SetExpiry(t time.Time) time.Time
	WriteToResponse(Response)
	GetId() string
	Clear()
}
