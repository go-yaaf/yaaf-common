package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-yaaf/yaaf-common/entity"
	"net/http"
	"strings"
	"time"
)

// region HTTP Mail Message --------------------------------------------------------------------------------------------

// HTTP mail message implementation
type httpMailMessage struct {
	client      *httpMailClient
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
func (m *httpMailMessage) From(from string) IMailMessage {
	m.from = from
	return m
}

// To Set recipients mail addresses
func (m *httpMailMessage) To(to []string) IMailMessage {
	m.to = to
	return m
}

// Cc Set cc list mail addresses
func (m *httpMailMessage) Cc(cc []string) IMailMessage {
	m.cc = cc
	return m
}

// Bcc Set Bcc list mail addresses
func (m *httpMailMessage) Bcc(bcc []string) IMailMessage {
	m.bcc = bcc
	return m
}

// Subject Set subject
func (m *httpMailMessage) Subject(subject string) IMailMessage {
	m.subject = subject
	return m
}

// Body Set subject
func (m *httpMailMessage) Body(body string) IMailMessage {
	m.body = body
	return m
}

// HtmlBody Set HTML Body
func (m *httpMailMessage) HtmlBody(html string) IMailMessage {
	m.html = html
	return m
}

// Attachments set list of message attachments
func (m *httpMailMessage) Attachments(attachments []MailMessageAttachment) IMailMessage {
	m.attachments = attachments
	return m
}

// AddTo adds mail address to the To recipients list
func (m *httpMailMessage) AddTo(to ...string) IMailMessage {
	m.to = append(m.to, to...)
	return m
}

// AddCc adds mail address to the CC list
func (m *httpMailMessage) AddCc(cc ...string) IMailMessage {
	m.cc = append(m.cc, cc...)
	return m
}

// AddBcc adds mail address to the BCC list
func (m *httpMailMessage) AddBcc(bcc ...string) IMailMessage {
	m.bcc = append(m.bcc, bcc...)
	return m
}

// AddAttachments add attachment to the list
func (m *httpMailMessage) AddAttachments(attachments ...MailMessageAttachment) IMailMessage {
	m.attachments = append(m.attachments, attachments...)
	return m
}

// Attach add list of file paths as attachments
func (m *httpMailMessage) Attach(paths ...string) IMailMessage {
	for _, path := range paths {
		if att, err := getAttachment(path); err == nil {
			m.AddAttachments(att)
		}
	}
	return m
}

// Send mail message
func (m *httpMailMessage) Send() error {

	buffer, jer := m.buildMessage()
	if jer != nil {
		return jer
	}
	req, err := http.NewRequest("POST", m.client.uri, bytes.NewBuffer(buffer))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	for k, v := range m.client.headers {
		req.Header.Set(k, v)
	}

	// Create http client
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		return fmt.Errorf("server return status: %s", resp.Status)
	}
}

// Build mail message
func (m *httpMailMessage) buildMessage() ([]byte, error) {
	message := make(entity.Json)
	message["from"] = m.from
	message["to"] = m.to
	message["cc"] = m.cc
	message["bcc"] = m.bcc
	message["subject"] = m.subject

	if len(m.body) > 0 {
		message["body"] = m.body
	} else if len(m.html) > 0 {
		message["body"] = m.html
	}

	message["attachments"] = m.attachments
	return json.Marshal(message)
}

// endregion

// region HTTP Mail Client ---------------------------------------------------------------------------------------------

// HTTP mail client implementation
type httpMailClient struct {
	uri     string
	user    string
	headers map[string]string
}

// MailUsr returns mail server authentication user
func (c *httpMailClient) MailUsr() string {
	return c.user
}

// CreateTextMessage Create plain text message
func (c *httpMailClient) CreateTextMessage() IMailMessage {
	return &httpMailMessage{
		client: c,
		to:     make([]string, 0),
		cc:     make([]string, 0),
		bcc:    make([]string, 0),
		mime:   "text/plain",
	}
}

// CreateHtmlMessage Create HTML message
func (c *httpMailClient) CreateHtmlMessage() IMailMessage {
	return &httpMailMessage{
		client:      c,
		to:          make([]string, 0),
		cc:          make([]string, 0),
		bcc:         make([]string, 0),
		attachments: make([]MailMessageAttachment, 0),
		mime:        "text/html",
	}
}

// CreateJsonMessage Create Json message
func (c *httpMailClient) CreateJsonMessage() IMailMessage {
	return &httpMailMessage{
		client:      c,
		to:          make([]string, 0),
		cc:          make([]string, 0),
		bcc:         make([]string, 0),
		attachments: make([]MailMessageAttachment, 0),
		mime:        "application/json",
	}
}

// CreateTemplateMessage Create Template message
func (c *httpMailClient) CreateTemplateMessage(template TemplateName, variables map[string]string) IMailMessage {
	return &httpMailMessage{
		client:      c,
		to:          make([]string, 0),
		cc:          make([]string, 0),
		bcc:         make([]string, 0),
		attachments: make([]MailMessageAttachment, 0),
		template:    template,
		variables:   variables,
	}
}

// Build mail message
func (c *httpMailClient) buildMessage(m *httpMailMessage) string {
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

func newHttpMailClient(uri string, user string, headers map[string]string) IMailClient {
	return &httpMailClient{
		uri:     uri,
		user:    user,
		headers: headers,
	}
}
