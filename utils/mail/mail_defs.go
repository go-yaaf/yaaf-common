package mail

import (
	"fmt"
	"net/url"
	"strings"
)

type TemplateName string

// MailConfig Configure mail client parameters
type MailConfig struct {

	// Mail relay URI (type://host:port)
	MailRelayUri string

	// Mail Relay User
	MailRelayUser string

	// Mail Relay Password
	MailRelayPassword string

	// Flag to use TLS encrypted connection
	UseTLS bool
}

// IMailClient Mail client interface
type IMailClient interface {
	MailUsr() string
	CreateTextMessage() IMailMessage
	CreateHtmlMessage() IMailMessage
	CreateJsonMessage() IMailMessage
	CreateTemplateMessage(template TemplateName, variables map[string]string) IMailMessage
}

// MailMessageAttachment represents message attachment
type MailMessageAttachment struct {
	// The full file path to attach
	FileName string

	// MIME type, ignore this field, it will be set automatically
	ContentType string

	// Base64 content of the file (ignore this field)
	Base64Content string
}

// IMailMessage Mail message interface
type IMailMessage interface {
	From(from string) IMailMessage
	To(to []string) IMailMessage
	Cc(cc []string) IMailMessage
	Bcc(bcc []string) IMailMessage
	Subject(subject string) IMailMessage
	Body(body string) IMailMessage
	HtmlBody(html string) IMailMessage
	Attachments(attachments []MailMessageAttachment) IMailMessage
	Send() error
}

// NewMailClient is a Mail client factory method
func NewMailClient(config MailConfig) (IMailClient, error) {

	if config.MailRelayUser == "" {
		return nil, fmt.Errorf("empty user name for mail provider URI: %s", config.MailRelayUri)
	}

	uri, err := url.Parse(config.MailRelayUri)
	if err != nil {
		return nil, fmt.Errorf("mail provider URI: %s parsing failed: %s", config.MailRelayUri, err.Error())
	}

	scheme := strings.ToLower(uri.Scheme)

	if scheme == "smtp" {
		return newSmtpMailClient(uri.Hostname(), uri.Port(), config.MailRelayUser, config.MailRelayPassword, config.UseTLS), nil
	} else {
		return nil, fmt.Errorf("unsupported mail type: %s", scheme)
	}
}

type mailAddress struct {
	Email string
	Name  string
}

// Parse mail address from the format: name <user@mail.com>
func parseMailAddress(addr string) mailAddress {

	if strings.Contains(addr, "<") {
		parts := strings.Split(addr, "<")
		name := strings.TrimSpace(parts[0])
		email := strings.ReplaceAll(parts[1], ">", "")
		email = strings.TrimSpace(email)
		return mailAddress{Email: email, Name: name}
	} else {
		return mailAddress{Email: addr, Name: ""}
	}
}
