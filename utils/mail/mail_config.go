package mail

import (
	"fmt"
	"net/url"
	"strings"
)

// region MailConfig Fluent API ----------------------------------------------------------------------------------------

// MailConfig holds the configuration for creating a mail client.
// It uses a fluent API to set the configuration parameters.
type MailConfig struct {
	// MailRelayUri is the URI of the mail relay, e.g., "smtp://smtp.example.com:587" or "https://api.mailgun.net/v3".
	MailRelayUri string

	// MailRelayUser is the username for authentication with the mail relay.
	MailRelayUser string

	// MailRelayPassword is the password for authentication with the mail relay.
	MailRelayPassword string

	// UseTLS specifies whether to use TLS for the connection. This is primarily for SMTP.
	UseTLS bool

	// Headers is a map of HTTP headers to be sent with the request, applicable for HTTP-based mail clients.
	Headers map[string]string
}

// NewMailConfig creates a new MailConfig with the given URI.
//
// Parameters:
//
//	uri: The URI of the mail relay.
//
// Returns:
//
//	A new MailConfig instance.
func NewMailConfig(uri string) *MailConfig {
	return &MailConfig{
		MailRelayUri: uri,
		Headers:      make(map[string]string),
	}
}

// WithUser sets the mail relay username.
//
// Parameters:
//
//	value: The username.
//
// Returns:
//
//	The MailConfig instance for chaining.
func (mc *MailConfig) WithUser(value string) *MailConfig {
	mc.MailRelayUser = value
	return mc
}

// WithPassword sets the mail relay password.
//
// Parameters:
//
//	value: The password.
//
// Returns:
//
//	The MailConfig instance for chaining.
func (mc *MailConfig) WithPassword(value string) *MailConfig {
	mc.MailRelayPassword = value
	return mc
}

// WithTLS sets the TLS flag.
//
// Parameters:
//
//	value: The TLS flag.
//
// Returns:
//
//	The MailConfig instance for chaining.
func (mc *MailConfig) WithTLS(value bool) *MailConfig {
	mc.UseTLS = value
	return mc
}

// SetHeader sets an HTTP header for HTTP-based mail clients.
//
// Parameters:
//
//	header: The name of the header.
//	value: The value of the header.
//
// Returns:
//
//	The MailConfig instance for chaining.
func (mc *MailConfig) SetHeader(header, value string) *MailConfig {
	mc.Headers[header] = value
	return mc
}

// CreateClient creates a mail client based on the configuration.
// It determines the client type (SMTP or HTTP) from the URI scheme.
//
// Returns:
//
//	An IMailClient instance.
//	An error if the URI is invalid or the scheme is unsupported.
func (mc *MailConfig) CreateClient() (IMailClient, error) {
	uri, err := url.Parse(mc.MailRelayUri)
	if err != nil {
		return nil, fmt.Errorf("invalid mail provider URI %s: %w", mc.MailRelayUri, err)
	}

	scheme := strings.ToLower(uri.Scheme)
	switch scheme {
	case "smtp":
		return newSmtpMailClient(uri.Hostname(), uri.Port(), mc.MailRelayUser, mc.MailRelayPassword, mc.UseTLS), nil
	case "http", "https":
		return newHttpMailClient(mc.MailRelayUri, mc.MailRelayUser, mc.Headers), nil
	default:
		return nil, fmt.Errorf("unsupported mail scheme: %s", scheme)
	}
}

// endregion
