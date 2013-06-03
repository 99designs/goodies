package test_mailer

import (
	"github.com/99designs/goodies/mailer"
	"log"
)

type TestMailer struct {
	sent []mailer.Email
}

func (t *TestMailer) Send(e mailer.Email) error {
	t.sent = append(t.sent, e)
	log.Printf("%+v\n", e)
	return nil
}

func NewTestMailer() mailer.Mailer {
	return &TestMailer{make([]mailer.Email, 0)}
}
