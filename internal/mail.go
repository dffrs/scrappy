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

func NewMail(message string, config *config) *mail {
	return &mail{
		from:     config.from,
		to:       config.to,
		password: config.password,
		host:     config.host,
		port:     config.port,
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
