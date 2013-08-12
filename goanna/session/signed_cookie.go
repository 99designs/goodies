// Package session implements a cookie session store for goanna
package session

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/99designs/goodies/goanna"
	"log"
	"net/http"
	"time"
)

const (
	EXPIRY_KEY = "SESSION_EXPIRY_TIME"
)

type CookieSessionStore struct {
	goanna.SessionStore
	goanna.CookieSigner
	CookieName      string
	DefaultDuration time.Duration
}

func NewCookieSessionStore(name, secret string, defaultDuration time.Duration) CookieSessionStore {
	return CookieSessionStore{
		CookieSigner:    goanna.NewCookieSigner(secret),
		CookieName:      name,
		DefaultDuration: defaultDuration,
	}
}

func (ss CookieSessionStore) getSessionData(request *http.Request) (*sessionData, error) {
	cookie, err := request.Cookie(ss.CookieName)
	if err != nil {
		return nil, err
	}
	return ss.decodeSessionData(cookie.Value)
}

func (ss CookieSessionStore) decodeSessionData(cv string) (*sessionData, error) {
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

func (ss CookieSessionStore) GetSession(request *goanna.Request) goanna.Session {
	session := SignedCookieSession{
		name:          ss.CookieName,
		CookieSigner:  ss.CookieSigner,
		defaultExpiry: time.Now().Add(ss.DefaultDuration),
		data:          &sessionData{},
	}
	data, err := ss.getSessionData(request.Request)
	if err != nil {
		session.data = &sessionData{
			Id:    generateSessionId(),
			Store: make(map[string]string),
		}
	} else {
		session.data = data
	}

	return session
}

type sessionData struct {
	Id    string
	Store map[string]string
}

type SignedCookieSession struct {
	name string
	goanna.CookieSigner
	defaultExpiry time.Time
	data          *sessionData
}

func (s SignedCookieSession) GetId() string {
	return s.data.Id
}
func (s SignedCookieSession) String() string {
	return fmt.Sprintf("%+v", s.data)
}

func (s SignedCookieSession) Get(key string) string {
	return s.data.Store[key]
}

func (s SignedCookieSession) Set(key string, value string) {
	s.data.Store[key] = value
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
	bytes := marshalSessionData(*s.data)

	cookie := http.Cookie{
		Name:     s.name,
		Value:    s.CookieSigner.EncodeRawData(bytes),
		Expires:  s.Expiry(),
		HttpOnly: true,
		Path:     "/",
	}
	resp.SetCookie(cookie)
}

func (s SignedCookieSession) Clear() {
	s.data.Id = generateSessionId()
	s.data.Store = make(map[string]string)
}

func generateSessionId() string {
	return randString(22) // aprox 128 bits of entropy (62^22)
}

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
