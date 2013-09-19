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

func NopSessionFinder(_ *Request) Session { return &NopSession{} }

type NopSession struct{}

func (s NopSession) Get(_ string) string             { return "" }
func (s NopSession) Set(_string, _ string)           {}
func (s NopSession) Expiry() time.Time               { return time.Now() }
func (s NopSession) SetExpiry(_ time.Time) time.Time { return time.Now() }
func (s NopSession) WriteToResponse(_ Response)      {}
func (s NopSession) GetId() string                   { return "" }
func (s NopSession) Clear()                          {}
