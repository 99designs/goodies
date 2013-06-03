// Package mailer provides email delivery (current delivery methods are 'postmark' and 'test')
package mailer

type Email struct {
	From    string
	ReplyTo string

	To  []string
	Cc  []string
	Bcc []string

	Subject  string
	HtmlBody string
	TextBody string

	Headers     map[string]string
	Attachments []Attachment
}

type Attachment struct {
	Name        string
	ContentType string
	Content     []byte
}

type Mailer interface {
	Send(Email) error
}

func NewEmail() Email {
	return Email{
		To:          make([]string, 0),
		Cc:          make([]string, 0),
		Bcc:         make([]string, 0),
		Attachments: make([]Attachment, 0),
		Headers:     make(map[string]string, 0),
	}
}
