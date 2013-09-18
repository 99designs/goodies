package session

import (
	"crypto/rand"
)

type sessionData struct {
	Id    string
	Store map[string]string
}

func (sd sessionData) GetId() string {
	return sd.Id
}

func (sd sessionData) Get(key string) string {
	return sd.Store[key]
}

func (sd sessionData) Set(key string, value string) {
	sd.Store[key] = value
}

func (sd sessionData) Clear() {
	sd.Id = generateSessionId()
	sd.Store = make(map[string]string)
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
