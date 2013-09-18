// Package session implements a cookie session store for goanna
package session

import (
	"bytes"
	"encoding/gob"
	"errors"
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
	return ss.decodeSessionData(cookie.Value)
}

func (ss SignedCookieSessionHandler) decodeSessionData(cv string) (*sessionData, error) {
	raw, err := ss.CookieSigner.DecodeCookieBytes(cv)
	if err != nil {
		return nil, err
	}
	return unmarshalSessionData(raw)
}

func marshalSessionData(sd sessionData) []byte {
	buf := &bytes.Buffer{}
	e := gob.NewEncoder(buf)
	err := e.Encode(sd)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func unmarshalSessionData(raw []byte) (*sessionData, error) {
	var data sessionData
	buf := bytes.NewBuffer(raw)
	d := gob.NewDecoder(buf)
	err := d.Decode(&data)
	if err != nil {
		log.Println("Invalid session cookie: " + err.Error())
		return nil, err
	}
	if data.Id == "" || data.Store == nil {
		return nil, errors.New("Nil data in struct")
	}
	return &data, nil
}

func (ss SignedCookieSessionHandler) GetSession(request *goanna.Request) goanna.Session {
	session := SignedCookieSession{
		SignedCookieSessionHandler: &ss,
	}
	data, err := ss.getSessionData(request.Request)
	if err == nil {
		session.sessionData = data
		if session.HasExpired() {
			session.Clear()
		}
	} else {
		session.sessionData = &sessionData{
			Id:    generateSessionId(),
			Store: make(map[string]string),
		}
	}

	return session
}

func (ss SignedCookieSessionHandler) writeToResponse(s SignedCookieSession, resp goanna.Response) {
	bytes := marshalSessionData(*s.sessionData)
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
	*SignedCookieSessionHandler
}

func (s SignedCookieSession) MaxAge() time.Duration {
	expiry, err := s.expiry()
	if err == nil {
		return expiry.Sub(time.Now())
	} else if s.SignedCookieSessionHandler.DefaultDuration > 0 {
		return s.SignedCookieSessionHandler.DefaultDuration
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
	s.SignedCookieSessionHandler.writeToResponse(s, resp)
}
