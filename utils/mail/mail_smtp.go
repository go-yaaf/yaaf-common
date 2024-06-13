package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

// region SMTP Mail Message --------------------------------------------------------------------------------------------

// SMTP mail message implementation
type smtpMailMessage struct {
	client      *smtpMailClient
	from        string
	to          []string
	cc          []string
	subject     string
	body        string
	html        string
	mime        string
	template    TemplateName
	attachments []MailMessageAttachment
	variables   map[string]string
}

// From Set sender mail address
func (m *smtpMailMessage) From(from string) IMailMessage {
	m.from = from
	return m
}

// To Set recipients mail addresses
func (m *smtpMailMessage) To(to []string) IMailMessage {
	m.to = to
	return m
}

// Cc Set cc list mail addresses
func (m *smtpMailMessage) Cc(cc []string) IMailMessage {
	m.cc = cc
	return m
}

// Subject Set subject
func (m *smtpMailMessage) Subject(subject string) IMailMessage {
	m.subject = subject
	return m
}

// Body Set subject
func (m *smtpMailMessage) Body(body string) IMailMessage {
	m.body = body
	return m
}

// HtmlBody Set HTML Body
func (m *smtpMailMessage) HtmlBody(html string) IMailMessage {
	m.html = html
	return m
}

// Attachments set list of message attachments
func (m *smtpMailMessage) Attachments(attachments []MailMessageAttachment) IMailMessage {
	m.attachments = attachments
	return m
}

// Send mail message
func (m *smtpMailMessage) Send() error {
	return m.client.send(m)
}

// endregion

// region SMTP Mail Client ---------------------------------------------------------------------------------------------

// SMTP mail client implementation
type smtpMailClient struct {
	host     string
	port     int
	user     string
	password string
	useTls   bool
}

// MailUsr set mail server authentication user
func (c *smtpMailClient) MailUsr() string {
	return c.user
}

// CreateTextMessage Create plain text message
func (c *smtpMailClient) CreateTextMessage() IMailMessage {
	return &smtpMailMessage{
		client: c,
		to:     make([]string, 0),
		cc:     make([]string, 0),
		mime:   "text/plain",
	}
}

// CreateHtmlMessage Create HTML message
func (c *smtpMailClient) CreateHtmlMessage() IMailMessage {
	return &smtpMailMessage{
		client: c,
		to:     make([]string, 0),
		cc:     make([]string, 0),
		mime:   "text/html",
	}
}

// CreateJsonMessage Create Json message
func (c *smtpMailClient) CreateJsonMessage() IMailMessage {
	return &smtpMailMessage{
		client: c,
		to:     make([]string, 0),
		cc:     make([]string, 0),
		mime:   "application/json",
	}
}

// CreateTemplateMessage Create Template message
func (c *smtpMailClient) CreateTemplateMessage(template TemplateName, variables map[string]string) IMailMessage {
	return &smtpMailMessage{
		client:    c,
		to:        make([]string, 0),
		cc:        make([]string, 0),
		template:  template,
		variables: variables,
	}
}

// Build mail message
func (c *smtpMailClient) buildMessage(m *smtpMailMessage) string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", m.from)
	if len(m.to) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(m.to, ";"))
	}
	if len(m.cc) > 0 {
		message += fmt.Sprintf("Cc: %s\r\n", strings.Join(m.cc, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", m.subject)
	message += "\r\n" + m.body

	return message
}

// endregion

func newSmtpMailClient(host, port, usr, pwd string, tls bool) IMailClient {

	var err error
	p := 80

	if p, err = strconv.Atoi(port); err != nil {
		p = 80
	}
	return &smtpMailClient{
		host:     host,
		port:     p,
		user:     usr,
		password: pwd,
		useTls:   tls,
	}
}

// Send mail
func (c *smtpMailClient) send(m *smtpMailMessage) (retError error) {

	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", m.to...)
	msg.SetHeader("Cc", m.cc...)
	msg.SetHeader("Subject", m.subject)

	if m.mime == "text/html" {
		msg.SetBody(m.mime, m.html)
	} else {
		msg.SetBody(m.mime, m.body)
	}

	f := func(w io.Writer) error {
		data, _ := base64.StdEncoding.DecodeString(m.attachments[0].Base64Content)
		_, err := io.Copy(w, bytes.NewReader(data))
		return err
	}
	if len(m.attachments) > 0 {
		msg.Attach("event_image.jpg", gomail.SetCopyFunc(f))
	}

	d := gomail.NewDialer(c.host, c.port, c.user, c.password)

	if err := d.DialAndSend(msg); err != nil {
		// test for specific error
		switch e := err.(type) {
		case *net.OpError:
			{
				netErr := err.(*net.OpError)
				switch netErr.Err.(type) {
				case *net.DNSError:
					retError = netErr.Err
				default:
					if netErr.Error() == "i/o timeout" {
						retError = netErr.Err
					}
				}
			}
		case *textproto.Error:
			{
				retError = err
			}
		default:
			retError = e
		}
	}
	return
}
