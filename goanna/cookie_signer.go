package goanna

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"net/http"
	"strings"
)

// CookieSigner signs a cookie with a sha236 hmac
// to ensure that it cannot be tampered with
type CookieSigner struct {
	signer hash.Hash
}

// NewCookieSigner returns a new CookieSigner using the given key
func NewCookieSigner(key []byte) CookieSigner {
	return CookieSigner{signer: hmac.New(sha256.New, key)}
}

// EncodeCookie signs and encodes the cookie
func (c CookieSigner) EncodeCookie(cookie *http.Cookie) {
	cookie.Value = c.EncodeValue(cookie.Value)
}

// DecodeCookie verifies the signature and decodes the value into the cookie
func (c CookieSigner) DecodeCookie(cookie *http.Cookie) error {
	data, err := c.DecodeValue(cookie.Value)
	if err != nil {
		return err
	}
	cookie.Value = data

	return nil
}

// DecodeValue validates and decodes a cookie value
func (c CookieSigner) DecodeValue(encodedvalue string) (string, error) {
	parts := strings.SplitN(encodedvalue, ".", 2)
	if len(parts) != 2 {
		return "", errors.New("Wrong number of parts")
	}
	mac, err := base64.URLEncoding.DecodeString(parts[0])
	value := parts[1]

	if err != nil {
		return "", err
	}
	if !hmac.Equal(c.mac(value), mac) {
		return "", errors.New("Bad signature")
	}

	return value, nil
}

// EncodeValue signs and encodes a cookie value
func (c CookieSigner) EncodeValue(value string) string {
	return fmt.Sprintf("%s.%s",
		base64.URLEncoding.EncodeToString(c.mac(value)),
		value)
}

func (c CookieSigner) mac(data string) []byte {
	c.signer.Reset()
	bytes := make([]byte, 0)
	fmt.Fprintf(c.signer, data)

	return c.signer.Sum(bytes)
}
