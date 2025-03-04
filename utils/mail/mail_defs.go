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

type TemplateName string

// region MailAddress Object --------------------------------------------------------------------------------------------

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

// endregion

// region IMailClient Interface ----------------------------------------------------------------------------------------

// IMailClient Mail client interface
type IMailClient interface {
	MailUsr() string
	CreateTextMessage() IMailMessage
	CreateHtmlMessage() IMailMessage
	CreateJsonMessage() IMailMessage
	CreateTemplateMessage(template TemplateName, variables map[string]string) IMailMessage
}

// endregion

// region Mail message attachment structure ----------------------------------------------------------------------------

// MailMessageAttachment represents message attachment
type MailMessageAttachment struct {
	// File name or the full file path (in case the file content is not provided) to attach
	FileName string `json:"fileName"`

	// Base64 content of the file
	Content string `json:"content"`

	// MIME type, if the content is provided as base64 in the Content field, it should include the suffix: ;base64
	ContentType string `json:"contentType"`
}

// endregion

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

	AddTo(to ...string) IMailMessage
	AddCc(cc ...string) IMailMessage
	AddBcc(bcc ...string) IMailMessage
	AddAttachments(attachments ...MailMessageAttachment) IMailMessage
	Attach(paths ...string) IMailMessage

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
	} else if scheme == "http" || scheme == "https" {
		return newHttpMailClient(config.MailRelayUri, config.MailRelayUser, config.Headers), nil
	} else {
		return nil, fmt.Errorf("unsupported mail type: %s", scheme)
	}
}

// Get Attachment from file path
func getAttachment(path string) (MailMessageAttachment, error) {
	result := MailMessageAttachment{}

	file, err := os.Open(path)
	if err != nil {
		return result, err
	}
	defer func() {
		_ = file.Close()
	}()

	content, er := io.ReadAll(file)
	if er != nil {
		return result, err
	}

	result.Content = base64.StdEncoding.EncodeToString(content)
	ext := filepath.Ext(path)
	contentType := mime.TypeByExtension(ext)
	if len(contentType) == 0 {
		contentType = "application/octet-stream"
	}
	result.ContentType = fmt.Sprintf("%s;base64", contentType)
	result.FileName = filepath.Base(path)
	return result, nil
}
