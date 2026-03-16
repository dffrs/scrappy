package internal

import (
	"errors"
	"fmt"
	"net/smtp"
)

type mail struct {
	from     string
	password string
	host     string
	port     int
	to       []string
	message  *string
}

func NewMail(config *config) *mail {
	return &mail{
		from:     config.from,
		to:       config.to,
		password: config.password,
		host:     config.host,
		port:     config.port,
	}
}

func (m *mail) SetMessage(message *string) *mail {
	m.message = message
	return m
}

func (m *mail) Send() error {
	if m.message == nil {
		return errors.New("message has not been set")
	}

	auth := smtp.PlainAuth("", m.from, m.password, m.host)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", m.host, m.port),
		auth,
		m.from,
		m.to,
		[]byte(*m.message),
	)
	if err != nil {
		return err
	}

	return nil
}
