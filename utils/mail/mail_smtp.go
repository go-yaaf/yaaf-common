package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
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
	bcc         []string
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

// Bcc Set Bcc list mail addresses
func (m *smtpMailMessage) Bcc(bcc []string) IMailMessage {
	m.bcc = bcc
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

// AddTo adds mail address to the To recipients list
func (m *smtpMailMessage) AddTo(to ...string) IMailMessage {
	m.to = append(m.to, to...)
	return m
}

// AddCc adds mail address to the CC list
func (m *smtpMailMessage) AddCc(cc ...string) IMailMessage {
	m.cc = append(m.cc, cc...)
	return m
}

// AddBcc adds mail address to the BCC list
func (m *smtpMailMessage) AddBcc(bcc ...string) IMailMessage {
	m.bcc = append(m.bcc, bcc...)
	return m
}

// AddAttachments add attachment to the list
func (m *smtpMailMessage) AddAttachments(attachments ...MailMessageAttachment) IMailMessage {
	m.attachments = append(m.attachments, attachments...)
	return m
}

// Attach add list of file paths as attachments
func (m *smtpMailMessage) Attach(paths ...string) IMailMessage {
	for _, path := range paths {
		if att, err := getAttachment(path); err == nil {
			m.AddAttachments(att)
		}
	}
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
		bcc:    make([]string, 0),
		mime:   "text/plain",
	}
}

// CreateHtmlMessage Create HTML message
func (c *smtpMailClient) CreateHtmlMessage() IMailMessage {
	return &smtpMailMessage{
		client: c,
		to:     make([]string, 0),
		cc:     make([]string, 0),
		bcc:    make([]string, 0),
		mime:   "text/html",
	}
}

// CreateJsonMessage Create Json message
func (c *smtpMailClient) CreateJsonMessage() IMailMessage {
	return &smtpMailMessage{
		client: c,
		to:     make([]string, 0),
		cc:     make([]string, 0),
		bcc:    make([]string, 0),
		mime:   "application/json",
	}
}

// CreateTemplateMessage Create Template message
func (c *smtpMailClient) CreateTemplateMessage(template TemplateName, variables map[string]string) IMailMessage {
	return &smtpMailMessage{
		client:    c,
		to:        make([]string, 0),
		cc:        make([]string, 0),
		bcc:       make([]string, 0),
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
	if len(m.bcc) > 0 {
		message += fmt.Sprintf("Bcc: %s\r\n", strings.Join(m.bcc, ";"))
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
	msg.SetHeader("Bcc", m.bcc...)
	msg.SetHeader("Subject", m.subject)

	if m.mime == "text/html" {
		if len(m.html) > 0 {
			msg.SetBody(m.mime, m.html)
		} else {
			msg.SetBody(m.mime, m.body)
		}
	} else {
		if len(m.body) > 0 {
			msg.SetBody(m.mime, m.body)
		} else {
			msg.SetBody(m.mime, m.html)
		}
	}

	// Complete attachments info
	for _, attachment := range m.attachments {
		if len(attachment.Content) == 0 {
			msg.Attach(attachment.FileName)
		} else {
			f := func(w io.Writer) error {
				data, _ := base64.StdEncoding.DecodeString(attachment.Content)
				_, err := io.Copy(w, bytes.NewReader(data))
				return err
			}
			fileName := filepath.Base(attachment.FileName)
			msg.Attach(fileName, gomail.SetCopyFunc(f))
		}
	}

	//f := func(w io.Writer) error {
	//	data, _ := base64.StdEncoding.DecodeString(m.attachments[0].Base64Content)
	//	_, err := io.Copy(w, bytes.NewReader(data))
	//	return err
	//}
	//if len(m.attachments) > 0 {
	//	msg.Attach("event_image.jpg", gomail.SetCopyFunc(f))
	//}

	d := gomail.NewDialer(c.host, c.port, c.user, c.password)
	return d.DialAndSend(msg)

	/*
		if err := d.DialAndSend(msg); err != nil {
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
	*/
}

func (c *smtpMailClient) attach(msg *gomail.Message, att MailMessageAttachment) (retError error) {

	if len(att.Content) == 0 {
		if err := c.attachPath(&att); err != nil {
			return err
		}
	}

	f := func(w io.Writer) error {
		data, _ := base64.StdEncoding.DecodeString(att.Content)
		_, err := io.Copy(w, bytes.NewReader(data))
		return err
	}
	msg.Attach(att.FileName, gomail.SetCopyFunc(f))
	return nil
}

func (c *smtpMailClient) attachPath(att *MailMessageAttachment) error {

	file, err := os.Open(att.FileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	content, er := io.ReadAll(file)
	if er != nil {
		return err
	}
	att.Content = base64.StdEncoding.EncodeToString(content)
	if len(att.ContentType) == 0 {
		ext := filepath.Ext(att.FileName)
		att.ContentType = mime.TypeByExtension(ext)
		if len(att.ContentType) == 0 {
			att.ContentType = "application/octet-stream"
		}
	}
	return nil
}
