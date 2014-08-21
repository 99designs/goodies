package goanna

import (
	"time"
)

// SessionFinder finds a session based on the request
type SessionFinder func(*Request) Session

type GoannaSessionHandlerFunc func(r *Request, s Session) Response

func (finder SessionFinder) Handler(handler GoannaSessionHandlerFunc) GoannaHandlerFunc {
	return func(r *Request) Response {
		sess := finder(r)
		resp := handler(r, sess)
		sess.WriteToResponse(resp)
		return resp
	}
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

func NopSessionFinder(_ *Request) Session { return &NopSession{} }

type NopSession struct{}

func (s NopSession) GetId() string              { return "" }
func (s NopSession) Get(_ string) string        { return "" }
func (s NopSession) Set(_string, _ string)      {}
func (s NopSession) SetMaxAge(_ time.Duration)  {}
func (s NopSession) HasExpired() bool           { return false }
func (s NopSession) Clear()                     {}
func (s NopSession) WriteToResponse(_ Response) {}
