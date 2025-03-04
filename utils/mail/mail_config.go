package mail

import (
	"fmt"
	"net/url"
	"strings"
)

// region MailConfig Fluent API ----------------------------------------------------------------------------------------

// MailConfig Configure mail client parameters
type MailConfig struct {

	// Mail relay URI (type://host:port)
	MailRelayUri string

	// Mail Relay User
	MailRelayUser string

	// Mail Relay Password
	MailRelayPassword string

	// Flag to use TLS encrypted connection (applicable for SMTP mail only)
	UseTLS bool

	// List of HTTP headers (applicable for HTTP mail only)
	Headers map[string]string
}

// NewMailConfig factory method
func NewMailConfig(uri string) *MailConfig {
	return &MailConfig{
		MailRelayUri: uri,
		Headers:      make(map[string]string),
	}
}

// WithUser set mail relay user
func (mc *MailConfig) WithUser(value string) *MailConfig {
	mc.MailRelayUser = value
	return mc
}

// WithPassword set mail relay password
func (mc *MailConfig) WithPassword(value string) *MailConfig {
	mc.MailRelayPassword = value
	return mc
}

// WithTLS set TLS flag
func (mc *MailConfig) WithTLS(value bool) *MailConfig {
	mc.UseTLS = value
	return mc
}

// SetHeader set HTTP header
func (mc *MailConfig) SetHeader(header, value string) *MailConfig {
	mc.Headers[header] = value
	return mc
}

// CreateClient create concrete implementation of mail client
func (mc *MailConfig) CreateClient() (IMailClient, error) {

	uri, err := url.Parse(mc.MailRelayUri)
	if err != nil {
		return nil, fmt.Errorf("mail provider URI: %s parsing failed: %s", mc.MailRelayUri, err.Error())
	}
	scheme := strings.ToLower(uri.Scheme)

	// Check for SMTP
	if scheme == "smtp" {
		return newSmtpMailClient(uri.Hostname(), uri.Port(), mc.MailRelayUser, mc.MailRelayPassword, mc.UseTLS), nil
	}

	// Check for HTTP(S)
	if scheme == "http" || scheme == "https" {
		return newHttpMailClient(mc.MailRelayUri, mc.MailRelayUser, mc.Headers), nil
	}

	// At this point, the schema is not supported
	return nil, fmt.Errorf("unsupported mail type: %s", scheme)
}

// endregion
