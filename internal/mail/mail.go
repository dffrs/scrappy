package mail

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"

	"scrappy/internal/types"
)

//go:embed template.html
var templateHTML string

type Mail struct {
	from     string
	password string
	host     string
	port     int
	to       []string
	subject  string
	products []types.ProductChanged
}

func NewMail() (*Mail, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	return &Mail{
		from:     config.From,
		to:       config.To,
		password: config.Password,
		host:     config.Host,
		port:     config.Port,
		subject:  "",
		products: nil,
	}, nil
}

func (m *Mail) SetSubject(subject string) *Mail {
	m.subject = subject
	return m
}

func (m *Mail) SetProducts(products []types.ProductChanged) *Mail {
	m.products = products
	return m
}

func (m *Mail) buildHeader() string {
	return "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
}

func (m *Mail) Send() error {
	if m.subject == "" {
		return errors.New("subject has not been set")
	}

	if m.products == nil {
		return errors.New("products have not been set")
	}

	var body bytes.Buffer
	t, err := template.ParseFiles(templateHTML)
	if err != nil {
		return err
	}

	err = t.Execute(&body, struct{ Products []types.ProductChanged }{Products: m.products})
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
