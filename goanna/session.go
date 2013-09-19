package goanna

import (
	"time"
)

// SessionFinder finds a session based on the request
type SessionFinder func(*Request) Session

type Session interface {
	GetId() string
	Get(string) string
	Set(string, string)
	SetMaxAge(time.Duration)
	HasExpired() bool
	Clear()
	WriteToResponse(Response)
}
