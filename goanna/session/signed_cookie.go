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

type CookieSessionHandler struct {
	goanna.SessionHandler
	goanna.CookieSigner
	CookieName      string
	DefaultDuration time.Duration
}

func NewCookieSessionHandler(name, secret string, defaultDuration time.Duration) CookieSessionHandler {
	return CookieSessionHandler{
		CookieSigner:    goanna.NewCookieSigner(secret),
		CookieName:      name,
		DefaultDuration: defaultDuration,
	}
}

func (ss CookieSessionHandler) getSessionData(request *http.Request) (*sessionData, error) {
	cookie, err := request.Cookie(ss.CookieName)
	if err != nil {
		return nil, err
	}
	return ss.decodeSessionData(cookie.Value)
}

func (ss CookieSessionHandler) decodeSessionData(cv string) (*sessionData, error) {
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

func (ss CookieSessionHandler) GetSession(request *goanna.Request) goanna.Session {
	session := SignedCookieSession{
		sessionData:   &sessionData{},
		name:          ss.CookieName,
		CookieSigner:  ss.CookieSigner,
		defaultExpiry: time.Now().Add(ss.DefaultDuration),
	}
	data, err := ss.getSessionData(request.Request)
	if err != nil {
		session.sessionData = &sessionData{
			Id:    generateSessionId(),
			Store: make(map[string]string),
		}
	} else {
		session.sessionData = data
	}

	return session
}

type SignedCookieSession struct {
	*sessionData
	name string
	goanna.CookieSigner
	defaultExpiry time.Time
}

func (s SignedCookieSession) SetDefaultExpiry() time.Time {
	s.SetExpiry(s.defaultExpiry)
	return s.defaultExpiry
}

func (s SignedCookieSession) Expiry() time.Time {
	expiryStr := s.Get(EXPIRY_KEY)
	expiry, err := time.Parse(time.RFC3339, expiryStr)
	if err != nil {
		return s.SetDefaultExpiry()
	}
	if expiry.After(s.defaultExpiry) {
		return expiry
	}
	return s.SetDefaultExpiry()
}

func (s SignedCookieSession) SetExpiry(e time.Time) time.Time {
	s.Set(EXPIRY_KEY, e.Format(time.RFC3339))
	return e
}

func (s SignedCookieSession) WriteToResponse(resp goanna.Response) {
	bytes := marshalSessionData(*s.sessionData)

	cookie := http.Cookie{
		Name:     s.name,
		Value:    s.CookieSigner.EncodeRawData(bytes),
		Expires:  s.Expiry(),
		HttpOnly: true,
		Path:     "/",
	}
	resp.SetCookie(cookie)
}
