// Package session implements a cookie session store for goanna
package session

import (
	"encoding/base64"
	"github.com/99designs/goodies/goanna"
	"log"
	"net/http"
	"time"
)

const (
	EXPIRY_KEY = "SESSION_EXPIRY_TIME"
)

type signedCookieSessionHandler struct {
	goanna.CookieSigner
	CookieName      string
	DefaultDuration time.Duration
	Secure          bool
}

func NewSignedCookieSessionFinder(name string, key []byte, defaultDuration time.Duration, secure bool) goanna.SessionFinder {
	ss := signedCookieSessionHandler{
		CookieSigner:    goanna.NewCookieSigner(key),
		CookieName:      name,
		DefaultDuration: defaultDuration,
		Secure:          secure,
	}

	return func(r *goanna.Request) goanna.Session {
		return ss.GetSession(r)
	}
}

func (ss signedCookieSessionHandler) getSessionData(request *http.Request) (*sessionData, error) {
	cookie, err := request.Cookie(ss.CookieName)
	if err != nil {
		return nil, err
	}

	err = ss.CookieSigner.DecodeCookie(cookie)
	if err != nil {
		return nil, err
	}

	decodedBytes, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, err
	}

	sessionData := sessionData{}
	err = sessionData.Marshal(decodedBytes)
	if err != nil {
		return nil, err
	}

	return &sessionData, nil
}

func (ss signedCookieSessionHandler) GetSession(request *goanna.Request) SignedCookieSession {
	session := SignedCookieSession{h: &ss}

	// Internet Explorer 8 and below do not support MaxAge
	ieVersion, err := getInternetExplorerVersion(request.UserAgent())
	if err == nil && ieVersion < 9 {
		session.preferExpires = true
	}

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

func (ss signedCookieSessionHandler) writeToResponse(s SignedCookieSession, resp goanna.Response) {
	bytes, err := s.sessionData.Unmarshal()
	if err != nil {
		log.Println(err.Error())
	}

	cookie := http.Cookie{
		Name:     ss.CookieName,
		Value:    base64.URLEncoding.EncodeToString(bytes),
		HttpOnly: true,
		Secure:   ss.Secure,
		Path:     "/",
	}
	ss.CookieSigner.EncodeCookie(&cookie)

	maxage := int(s.MaxAge() / time.Second)
	if maxage != 0 {
		if s.preferExpires {
			cookie.Expires = time.Now().Add(s.MaxAge())
		} else {
			cookie.MaxAge = maxage
		}
	}
	resp.SetCookie(cookie)
}

type SignedCookieSession struct {
	*sessionData
	h             *signedCookieSessionHandler
	preferExpires bool
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
