package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-yaaf/yaaf-common/entity"
)

// region HTTP Mail Message --------------------------------------------------------------------------------------------

// httpMailMessage is an implementation of the IMailMessage interface for sending email via an HTTP API.
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

// From sets the sender's email address.
func (m *httpMailMessage) From(from string) IMailMessage {
	m.from = from
	return m
}

// To sets the recipients' email addresses.
func (m *httpMailMessage) To(to []string) IMailMessage {
	m.to = to
	return m
}

// Cc sets the CC recipients' email addresses.
func (m *httpMailMessage) Cc(cc []string) IMailMessage {
	m.cc = cc
	return m
}

// Bcc sets the BCC recipients' email addresses.
func (m *httpMailMessage) Bcc(bcc []string) IMailMessage {
	m.bcc = bcc
	return m
}

// Subject sets the email's subject line.
func (m *httpMailMessage) Subject(subject string) IMailMessage {
	m.subject = subject
	return m
}

// Body sets the plain text body of the email.
func (m *httpMailMessage) Body(body string) IMailMessage {
	m.body = body
	return m
}

// HtmlBody sets the HTML body of the email.
func (m *httpMailMessage) HtmlBody(html string) IMailMessage {
	m.html = html
	return m
}

// Attachments sets the list of attachments for the email.
func (m *httpMailMessage) Attachments(attachments []MailMessageAttachment) IMailMessage {
	m.attachments = attachments
	return m
}

// AddTo adds one or more recipients to the To list.
func (m *httpMailMessage) AddTo(to ...string) IMailMessage {
	m.to = append(m.to, to...)
	return m
}

// AddCc adds one or more recipients to the CC list.
func (m *httpMailMessage) AddCc(cc ...string) IMailMessage {
	m.cc = append(m.cc, cc...)
	return m
}

// AddBcc adds one or more recipients to the BCC list.
func (m *httpMailMessage) AddBcc(bcc ...string) IMailMessage {
	m.bcc = append(m.bcc, bcc...)
	return m
}

// AddAttachments adds one or more attachments to the email.
func (m *httpMailMessage) AddAttachments(attachments ...MailMessageAttachment) IMailMessage {
	m.attachments = append(m.attachments, attachments...)
	return m
}

// Attach adds attachments from a list of file paths.
func (m *httpMailMessage) Attach(paths ...string) IMailMessage {
	for _, path := range paths {
		if att, err := getAttachment(path); err == nil {
			m.AddAttachments(att)
		}
	}
	return m
}

// Send sends the email using the HTTP client.
// It builds the message payload and sends it as a POST request.
func (m *httpMailMessage) Send() error {
	payload, err := m.buildMessage()
	if err != nil {
		return fmt.Errorf("failed to build http message: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, m.client.uri, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	for k, v := range m.client.headers {
		req.Header.Set(k, v)
	}

	// Create and send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http request failed with status: %s", resp.Status)
	}
	return nil
}

// buildMessage constructs the JSON payload for the HTTP API call.
func (m *httpMailMessage) buildMessage() ([]byte, error) {
	message := make(entity.Json)
	message["from"] = m.from
	message["to"] = m.to
	message["cc"] = m.cc
	message["bcc"] = m.bcc
	message["subject"] = m.subject

	if m.html != "" {
		message["body"] = m.html
	} else {
		message["body"] = m.body
	}

	if len(m.attachments) > 0 {
		message["attachments"] = m.attachments
	}

	return json.Marshal(message)
}

// endregion

// region HTTP Mail Client ---------------------------------------------------------------------------------------------

// httpMailClient is an implementation of the IMailClient interface that sends email via an HTTP API.
type httpMailClient struct {
	uri     string
	user    string
	headers map[string]string
}

// newHttpMailClient creates a new httpMailClient.
func newHttpMailClient(uri, user string, headers map[string]string) IMailClient {
	return &httpMailClient{
		uri:     uri,
		user:    user,
		headers: headers,
	}
}

// MailUsr returns the username used for authentication.
func (c *httpMailClient) MailUsr() string {
	return c.user
}

// CreateTextMessage creates a new plain text message.
func (c *httpMailClient) CreateTextMessage() IMailMessage {
	return &httpMailMessage{
		client: c,
		mime:   "text/plain",
	}
}

// CreateHtmlMessage creates a new HTML message.
func (c *httpMailClient) CreateHtmlMessage() IMailMessage {
	return &httpMailMessage{
		client: c,
		mime:   "text/html",
	}
}

// CreateJsonMessage creates a new JSON message.
func (c *httpMailClient) CreateJsonMessage() IMailMessage {
	return &httpMailMessage{
		client: c,
		mime:   "application/json",
	}
}

// CreateTemplateMessage creates a new message from a template.
func (c *httpMailClient) CreateTemplateMessage(template TemplateName, variables map[string]string) IMailMessage {
	return &httpMailMessage{
		client:    c,
		template:  template,
		variables: variables,
	}
}

// endregion
