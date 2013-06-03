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

type CookieSigner struct {
	signer hash.Hash
}

func NewCookieSigner(cookieSecret string) CookieSigner {
	return CookieSigner{signer: hmac.New(sha256.New, []byte(cookieSecret))}
}

func (c CookieSigner) EncodeCookie(cookie http.Cookie) *http.Cookie {
	cookie.Value = c.encodeValue(cookie.Value)

	return &cookie
}

func (c CookieSigner) DecodeCookie(cookie http.Cookie) (*http.Cookie, error) {
	data, err := c.DecodeCookieBytes(cookie.Value)
	if err != nil {
		return nil, err
	}
	cookie.Value = string(data)
	return &cookie, nil
}

func (c CookieSigner) DecodeCookieBytes(cookieValue string) ([]byte, error) {
	parts := strings.Split(cookieValue, ".")
	if len(parts) != 2 {
		return nil, errors.New("More than 2 parts")
	}
	rawdata, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	if c.encodeValue(string(rawdata)) != cookieValue {
		return nil, errors.New("Bad signature")
	}
	return rawdata, nil
}

func (c CookieSigner) EncodeRawData(d []byte) string {
	c.signer.Reset()
	bytes := make([]byte, 0)
	c.signer.Write(d)

	return fmt.Sprintf("%s.%s",
		base64.URLEncoding.EncodeToString([]byte(d)),
		base64.URLEncoding.EncodeToString(c.signer.Sum(bytes)))
}

func (c CookieSigner) encodeValue(data string) string {
	c.signer.Reset()
	bytes := make([]byte, 0)
	fmt.Fprintf(c.signer, data)

	return fmt.Sprintf("%s.%s",
		base64.URLEncoding.EncodeToString([]byte(data)),
		base64.URLEncoding.EncodeToString(c.signer.Sum(bytes)))
}
