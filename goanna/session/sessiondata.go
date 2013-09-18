package session

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"log"
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

func (sd sessionData) GobEncode() ([]byte, error) {
	buf := &bytes.Buffer{}
	e := gob.NewEncoder(buf)
	err := e.Encode(sd)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (data *sessionData) GobDecode(raw []byte) error {
	buf := bytes.NewBuffer(raw)
	d := gob.NewDecoder(buf)
	err := d.Decode(&data)
	if err != nil {
		log.Println("Invalid session cookie: " + err.Error())
		return err
	}
	if data.Id == "" || data.Store == nil {
		return errors.New("Nil data in struct")
	}
	return nil
}
