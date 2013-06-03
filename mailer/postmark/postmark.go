// Package postmark implements a backend for 'mailer' via postmark
package postmark

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/99designs/goodies/mailer"
	"io"
	"net/http"
	"strings"
)

const (
	postmarkUri  = "https://api.postmarkapp.com/email"
	jsonMimeType = "application/json"
	apiKeyHeader = "X-Postmark-Server-Token"
)

type PostmarkMailer struct {
	api_key string
	Tag     string
}

func NewMailer(api_key string) mailer.Mailer {
	return &PostmarkMailer{api_key: api_key}
}

type postmarkEmail struct {
	From        string
	ReplyTo     string
	To          string
	Cc          string
	Bcc         string
	Subject     string
	HtmlBody    string
	TextBody    string
	Headers     []postmarkHeader
	Attachments []mailer.Attachment
	Tag         string
}

type postmarkHeader struct {
	Name  string
	Value string
}

func (t *PostmarkMailer) postmarkEmail(e mailer.Email) postmarkEmail {
	var p postmarkEmail
	p.From = e.From
	p.ReplyTo = e.ReplyTo
	p.To = strings.Join(e.To, " , ")
	p.Cc = strings.Join(e.Cc, " , ")
	p.Bcc = strings.Join(e.Bcc, " , ")
	p.Subject = e.Subject
	p.HtmlBody = e.HtmlBody
	p.TextBody = e.TextBody
	p.Headers = make([]postmarkHeader, len(e.Headers))
	i := 0
	for k, v := range e.Headers {
		p.Headers[i].Name = k
		p.Headers[i].Value = v
		i += 1
	}
	p.Attachments = e.Attachments
	p.Tag = t.Tag
	return p
}

func (t *PostmarkMailer) postBody(e mailer.Email) (io.Reader, error) {
	body, err := json.Marshal(t.postmarkEmail(e))
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))

	return bytes.NewBuffer(body), nil
}

func (t *PostmarkMailer) buildRequest(e mailer.Email) (req *http.Request, err error) {
	var reader io.Reader
	reader, err = t.postBody(e)
	req, err = http.NewRequest("POST", postmarkUri, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", jsonMimeType)
	req.Header.Set("Accept", jsonMimeType)
	req.Header.Set(apiKeyHeader, t.api_key)
	return req, nil
}

func (t *PostmarkMailer) Send(e mailer.Email) (err error) {
	var req *http.Request
	var resp *http.Response

	req, err = t.buildRequest(e)
	if err != nil {
		return
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var j map[string]interface{}
		body := make([]byte, resp.ContentLength)
		_, err = resp.Body.Read(body)
		if err != nil {
			fmt.Println("Error reading resp.Body", string(body), err)
			return
		}
		err = json.Unmarshal(body, &j)
		if err != nil {
			return
		}
		if j == nil {
			return errors.New("Empty JSON body from postmark")
		}
		return errors.New(fmt.Sprintf("%+v: %+v (%s)", j["ErrorCode"], j["Message"], body))
	}
	return
}
