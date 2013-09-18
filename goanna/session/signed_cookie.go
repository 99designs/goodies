// Package session implements a cookie session store for goanna
package session

import (
	"github.com/99designs/goodies/goanna"
	"log"
	"net/http"
	"time"
)

const (
	EXPIRY_KEY = "SESSION_EXPIRY_TIME"
)

type SignedCookieSessionHandler struct {
	goanna.SessionHandler
	goanna.CookieSigner
	CookieName      string
	DefaultDuration time.Duration
}

func NewSignedCookieSessionHandler(name, secret string, defaultDuration time.Duration) SignedCookieSessionHandler {
	return SignedCookieSessionHandler{
		CookieSigner:    goanna.NewCookieSigner(secret),
		CookieName:      name,
		DefaultDuration: defaultDuration,
	}
}

func (ss SignedCookieSessionHandler) getSessionData(request *http.Request) (*sessionData, error) {
	cookie, err := request.Cookie(ss.CookieName)
	if err != nil {
		return nil, err
	}

	raw, err := ss.CookieSigner.DecodeCookieBytes(cookie.Value)
	if err != nil {
		return nil, err
	}

	sessionData := sessionData{}
	err = sessionData.GobDecode(raw)
	if err != nil {
		return nil, err
	}

	return &sessionData, nil
}

func (ss SignedCookieSessionHandler) GetSession(request *goanna.Request) goanna.Session {
	session := SignedCookieSession{h: &ss}
	data, err := ss.getSessionData(request.Request)
	if err == nil {
		session.sessionData = data
	} else {
		session.sessionData = NewSessionData()
	}

	if session.HasExpired() {
		session.Clear()
	}

	return session
}

func (ss SignedCookieSessionHandler) writeToResponse(s SignedCookieSession, resp goanna.Response) {
	bytes, err := s.sessionData.GobEncode()
	if err != nil {
		log.Println(err.Error())
	}

	signedbytes := ss.CookieSigner.EncodeRawData(bytes)

	cookie := http.Cookie{
		Name:     ss.CookieName,
		Value:    signedbytes,
		HttpOnly: true,
		Path:     "/",
	}
	maxage := int(s.MaxAge() / time.Second)
	if maxage != 0 {
		cookie.MaxAge = maxage
	}
	resp.SetCookie(cookie)
}

type SignedCookieSession struct {
	*sessionData
	h *SignedCookieSessionHandler
}

func (s SignedCookieSession) MaxAge() time.Duration {
	expiry, err := s.expiry()
	if err == nil {
		return expiry.Sub(time.Now())
	} else if s.h.DefaultDuration > 0 {
		return s.h.DefaultDuration
	}

	return 0
}

func (s SignedCookieSession) SetMaxAge(d time.Duration) {
	expiry := time.Now().Add(d)
	s.Set(EXPIRY_KEY, expiry.Format(time.RFC3339))
}

func (s SignedCookieSession) expiry() (time.Time, error) {
	expiryStr := s.Get(EXPIRY_KEY)
	return time.Parse(time.RFC3339, expiryStr)
}

func (s SignedCookieSession) HasExpired() bool {
	expiry, err := s.expiry()
	return err == nil && expiry.Before(time.Now())
}

func (s SignedCookieSession) WriteToResponse(resp goanna.Response) {
	s.h.writeToResponse(s, resp)
}
