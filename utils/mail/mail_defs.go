package mail

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// TemplateName is a type alias for a string, representing the name of a mail template.
type TemplateName string

// region MailAddress Object --------------------------------------------------------------------------------------------

// mailAddress represents an email address with an optional name.
type mailAddress struct {
	Email string
	Name  string
}

// parseMailAddress parses a string in the format "Name <user@example.com>" into a mailAddress struct.
// If the address is not in this format, it's treated as a simple email address.
func parseMailAddress(addr string) mailAddress {
	if strings.Contains(addr, "<") {
		parts := strings.Split(addr, "<")
		name := strings.TrimSpace(parts[0])
		email := strings.TrimSuffix(strings.TrimSpace(parts[1]), ">")
		return mailAddress{Email: email, Name: name}
	}
	return mailAddress{Email: addr}
}

// endregion

// region IMailClient Interface ----------------------------------------------------------------------------------------

// IMailClient defines the interface for a mail client, which is responsible for creating and sending messages.
type IMailClient interface {
	// MailUsr returns the username of the mail client.
	MailUsr() string

	// CreateTextMessage creates a new plain text mail message.
	CreateTextMessage() IMailMessage

	// CreateHtmlMessage creates a new HTML mail message.
	CreateHtmlMessage() IMailMessage

	// CreateJsonMessage creates a new JSON mail message.
	CreateJsonMessage() IMailMessage

	// CreateTemplateMessage creates a new mail message from a template.
	CreateTemplateMessage(template TemplateName, variables map[string]string) IMailMessage
}

// endregion

// region Mail message attachment structure ----------------------------------------------------------------------------

// MailMessageAttachment represents a file attached to an email.
type MailMessageAttachment struct {
	// FileName is the name of the file to be attached.
	FileName string `json:"fileName"`

	// Content is the base64-encoded content of the file.
	Content string `json:"content"`

	// ContentType is the MIME type of the attachment. For base64 content, it should include ";base64".
	ContentType string `json:"contentType"`
}

// endregion

// IMailMessage defines the interface for a mail message, providing methods to build and send the message.
type IMailMessage interface {
	From(from string) IMailMessage
	To(to []string) IMailMessage
	Cc(cc []string) IMailMessage
	Bcc(bcc []string) IMailMessage
	Subject(subject string) IMailMessage
	Body(body string) IMailMessage
	HtmlBody(html string) IMailMessage
	Attachments(attachments []MailMessageAttachment) IMailMessage

	AddTo(to ...string) IMailMessage
	AddCc(cc ...string) IMailMessage
	AddBcc(bcc ...string) IMailMessage
	AddAttachments(attachments ...MailMessageAttachment) IMailMessage
	Attach(paths ...string) IMailMessage

	Send() error
}

// NewMailClient is a factory function that creates a new mail client based on the provided configuration.
func NewMailClient(config MailConfig) (IMailClient, error) {
	if config.MailRelayUser == "" {
		return nil, fmt.Errorf("mail user is required for mail provider URI: %s", config.MailRelayUri)
	}

	uri, err := url.Parse(config.MailRelayUri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mail provider URI %s: %w", config.MailRelayUri, err)
	}

	scheme := strings.ToLower(uri.Scheme)
	switch scheme {
	case "smtp":
		return newSmtpMailClient(uri.Hostname(), uri.Port(), config.MailRelayUser, config.MailRelayPassword, config.UseTLS), nil
	case "http", "https":
		return newHttpMailClient(config.MailRelayUri, config.MailRelayUser, config.Headers), nil
	default:
		return nil, fmt.Errorf("unsupported mail scheme: %s", scheme)
	}
}

// getAttachment creates a MailMessageAttachment from a file path.
// It reads the file, base64-encodes its content, and determines the MIME type.
func getAttachment(path string) (MailMessageAttachment, error) {
	var result MailMessageAttachment

	file, err := os.Open(path)
	if err != nil {
		return result, fmt.Errorf("failed to open attachment file %s: %w", path, err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return result, fmt.Errorf("failed to read attachment file %s: %w", path, err)
	}

	result.Content = base64.StdEncoding.EncodeToString(content)
	ext := filepath.Ext(path)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	result.ContentType = fmt.Sprintf("%s;base64", contentType)
	result.FileName = filepath.Base(path)
	return result, nil
}
