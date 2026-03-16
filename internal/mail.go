package internal

import (
	"fmt"
	"net/smtp"
)

type mail struct {
	from     string
	password string
	host     string
	port     int
	to       []string
	message  string
}

func NewMail(message string) *mail {
	return &mail{
		from:     "",               // TODO: Load from .env
		password: "",               // TODO: Load from .env
		host:     "smtp.gmail.com", // TODO: Load from .env
		port:     587,              // TODO: Load from .env
		to:       []string{""},     // TODO: Load from .env
		message:  message,
	}
}

func (m *mail) Send() error {
	auth := smtp.PlainAuth("", m.from, m.password, m.host)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", m.host, m.port),
		auth,
		m.from,
		m.to,
		[]byte(m.message),
	)
	if err != nil {
		return err
	}

	return nil
}
