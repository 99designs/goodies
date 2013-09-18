package goanna

import (
	"time"
)

type SessionHandler interface {
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

type NopSessionStore struct{}

func (s NopSessionStore) GetSession(_ *Request) Session { return &NopSession{} }

type NopSession struct{}

func (s NopSession) Get(_ string) string             { return "" }
func (s NopSession) Set(_string, _ string)           {}
func (s NopSession) Expiry() time.Time               { return time.Now() }
func (s NopSession) SetExpiry(_ time.Time) time.Time { return time.Now() }
func (s NopSession) WriteToResponse(_ Response)      {}
func (s NopSession) GetId() string                   { return "" }
func (s NopSession) Clear()                          {}
