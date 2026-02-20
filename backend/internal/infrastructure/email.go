package infrastructure

import (
	"fmt"
	"basekarya-backend/internal/config"
	"io"

	"gopkg.in/gomail.v2"
)

type EmailProvider struct {
	dialer *gomail.Dialer
	from   string
}

func NewEmailProvider(cfg *config.EmailConfig) *EmailProvider {
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	return &EmailProvider{
		dialer: d,
		from:   cfg.From,
	}
}

func (e *EmailProvider) SendWithAttachment(to, subject, htmlBody, fileName string, attachmentBytes []byte) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("HRIS System <%s>", e.from))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	m.Attach(fileName, gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(attachmentBytes)
		return err
	}))

	return e.dialer.DialAndSend(m)
}
