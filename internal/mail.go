package internal

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"
)

type mail struct {
	from     string
	password string
	host     string
	port     int
	to       []string
	subject  *string
	message  *string
}

func (m *mail) buildHeader() string {
	return "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
}

func NewMail() (*mail, error) {
	config, err := loadEnv()
	if err != nil {
		return nil, err
	}

	return &mail{
		from:     config.from,
		to:       config.to,
		password: config.password,
		host:     config.host,
		port:     config.port,
		subject:  nil,
		message:  nil,
	}, nil
}

func (m *mail) SetSubject(subject *string) *mail {
	m.subject = subject
	return m
}

func (m *mail) SetMessage(message *string) *mail {
	m.message = message
	return m
}

func (m *mail) Send() error {
	if m.subject == nil {
		return errors.New("subject has not been set")
	}

	if m.message == nil {
		return errors.New("message has not been set")
	}

	var body bytes.Buffer
	// TODO: this needs to be embeded
	t, err := template.ParseFiles("assets/template.html")
	if err != nil {
		return err
	}

	err = t.Execute(&body, struct{ Name string }{Name: "Daniel"})
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Subject: %s\n%s\n\n%s", *m.subject, m.buildHeader(), body.String())

	auth := smtp.PlainAuth("", m.from, m.password, m.host)

	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", m.host, m.port),
		auth,
		m.from,
		m.to,
		[]byte(msg),
	)
	if err != nil {
		return err
	}

	return nil
}
