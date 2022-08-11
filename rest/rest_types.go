// Copyright 2022. Motty Cohen
//
// package types
//
package rest

import (
	"net/http"
)

// IRestResponse is an interface that all rest response messages must comply
type IRestResponse interface {
	GetErrorCode() int
	GetErrorMessage() string
	SetError(err error)
	GetRawContent() []byte
	SetRawContent(content []byte)
	GetMimeType() string
	SetMimeType(mimeType string)
	GetHeaders() map[string]string
	SetHeader(key string, value string)
}

// RequestWithToken includes access to the JWT token
type RequestWithToken struct {
	*http.Request
	Token          string
	ResponseWriter http.ResponseWriter
}

// NewRequestWithToken is the factory method of RequestWithToken
func NewRequestWithToken(rw http.ResponseWriter, r *http.Request, token string) *RequestWithToken {
	return &RequestWithToken{
		Request:        r,
		Token:          token,
		ResponseWriter: rw,
	}
}

// RestHandler function signature
type RestHandler func(r *RequestWithToken) (IRestResponse, error)

// RestHandlerAdaptorFunc function signature
type RestHandlerAdaptorFunc func(entry RestEntry) http.HandlerFunc

// region REST endpoints -----------------------------------------------------------------------------------------------

type RestEntry struct {
	Path    string      // Rest method path
	Method  string      // HTTP method verb
	Handler RestHandler // Handler function (http.HandlerFunc)
}

type IRestEndpoint interface {
	Entries() []RestEntry
}

type RestEndpoint struct {
	entries []RestEntry
}

func (r RestEndpoint) Entries() []RestEntry { return r.entries }

// endregion

// region Static files endpoints ---------------------------------------------------------------------------------------

type StaticFilesEntry struct {
	Path   string // Static content path
	Folder string // Static content folder
}

type IStaticEndpoint interface {
	Entries() []StaticFilesEntry
}

type StaticEndpoint struct {
	entries []StaticFilesEntry
}

func (r StaticEndpoint) Entries() []StaticFilesEntry { return r.entries }

// endregion
