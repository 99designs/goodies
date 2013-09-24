package goanna

import (
	"bytes"
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

func (c CookieSigner) EncodeCookie(cookie *http.Cookie) {
	cookie.Value = c.EncodeValue(cookie.Value)
}

func (c CookieSigner) DecodeCookie(cookie *http.Cookie) error {
	data, err := c.DecodeValue(cookie.Value)
	if err != nil {
		return err
	}
	cookie.Value = data

	return nil
}

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
	if !bytes.Equal(c.mac(value), mac) {
		return "", errors.New("Bad signature")
	}

	return value, nil
}

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
