package postmark

import (
	"github.com/99designs/goodies/mailer"
	"testing"
)

func TestPostmark(t *testing.T) {
	email := mailer.NewEmail()
	email.From = "noreply@account.99designs.com"
	email.To = append(email.To, "daniel@heath.cc")
	email.Subject = "Postmark Test"
	email.HtmlBody = "<html><body><strong>Hello</strong>, dear Postmark user.</body></html>"
	email.TextBody = "Hello, dear Postmark user."

	m := NewMailer("POSTMARK_API_TEST")
	err := m.Send(email)
	if err != nil {
		t.Error(err)
	}
}
