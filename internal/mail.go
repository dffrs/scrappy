package internal

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"

	"scrappy/types"
)

type mail struct {
	from     string
	password string
	host     string
	port     int
	to       []string
	subject  string
	products []types.Product
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
		subject:  "",
		products: nil,
	}, nil
}

func (m *mail) SetSubject(subject string) *mail {
	m.subject = subject
	return m
}

func (m *mail) SetProducts(products []types.Product) *mail {
	m.products = products
	return m
}

func (m *mail) Send() error {
	if m.subject == "" {
		return errors.New("subject has not been set")
	}

	if m.products == nil {
		return errors.New("products have not been set")
	}

	var body bytes.Buffer
	// TODO: this needs to be embeded
	t, err := template.ParseFiles("assets/template.html")
	if err != nil {
		return err
	}

	err = t.Execute(&body, struct {
		Products []types.Product
	}{Products: m.products})
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Subject: %s\n%s\n\n%s", m.subject, m.buildHeader(), body.String())

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
