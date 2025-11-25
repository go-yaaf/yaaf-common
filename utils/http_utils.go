// Package utils provides a collection of utility functions, including helpers for making HTTP requests.
package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/go-yaaf/yaaf-common/entity"
	"net/http"
	"net/url"
	"time"
)

// HttpUtilsStruct provides a fluent interface for building and sending HTTP requests.
// It allows for setting the method, URL, headers, body, and timeout for a request.
type HttpUtilsStruct struct {
	method     string
	url        string
	body       string
	headers    map[string]string
	TimeoutSec int
}

// HttpUtils is a factory function that creates a new instance of HttpUtilsStruct.
// This is the entry point for using the HTTP utility.
func HttpUtils() *HttpUtilsStruct {
	return &HttpUtilsStruct{
		method:  "GET",
		headers: make(map[string]string),
	}
}

// New initializes a new request with a specified method and URL.
//
// Parameters:
//
//	method: The HTTP method (e.g., "GET", "POST", "PUT").
//	url: The URL for the request.
//
// Returns:
//
//	The HttpUtilsStruct instance for chaining.
func (u *HttpUtilsStruct) New(method, url string) *HttpUtilsStruct {
	u.method = method
	u.url = url
	return u
}

// WithHeader adds a single header to the request.
//
// Parameters:
//
//	key: The header name.
//	value: The header value.
//
// Returns:
//
//	The HttpUtilsStruct instance for chaining.
func (u *HttpUtilsStruct) WithHeader(key, value string) *HttpUtilsStruct {
	u.headers[key] = value
	return u
}

// WithHeaders adds multiple headers to the request from a map.
//
// Parameters:
//
//	headers: A map of header names to header values.
//
// Returns:
//
//	The HttpUtilsStruct instance for chaining.
func (u *HttpUtilsStruct) WithHeaders(headers map[string]string) *HttpUtilsStruct {
	for k, v := range headers {
		u.headers[k] = v
	}
	return u
}

// WithBody sets the body of the request.
//
// Parameters:
//
//	body: The request body as a string.
//
// Returns:
//
//	The HttpUtilsStruct instance for chaining.
func (u *HttpUtilsStruct) WithBody(body string) *HttpUtilsStruct {
	u.body = body
	return u
}

// WithTimeout sets the timeout in seconds for the request.
// If not set or set to 0, there is no timeout.
//
// Parameters:
//
//	timeout: The timeout duration in seconds.
//
// Returns:
//
//	The HttpUtilsStruct instance for chaining.
func (u *HttpUtilsStruct) WithTimeout(timeout int) *HttpUtilsStruct {
	u.TimeoutSec = timeout
	return u
}

// Send constructs and sends the HTTP request based on the configured parameters.
// It handles URL parsing, basic authentication from the URL, and header setup.
//
// IMPORTANT: The caller is responsible for closing the response body, e.g., `defer res.Body.Close()`.
//
// Returns:
//
//	A pointer to the http.Response and an error if one occurred.
//	An error is also returned for non-successful status codes (i.e., not 2xx or 3xx).
func (u *HttpUtilsStruct) Send() (*http.Response, error) {
	parsedUrl, err := url.Parse(u.url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Add Basic Authentication header if user info is present in the URL.
	if parsedUrl.User != nil {
		userPassword := parsedUrl.User.String()
		// Clear user info from the URL to avoid it being sent in the request line.
		parsedUrl.User = nil
		authHeader := authenticationHeader(userPassword)
		u.headers[authHeader.Key] = authHeader.Value
	}

	// Create the request.
	req, err := http.NewRequest(u.method, parsedUrl.String(), bytes.NewBufferString(u.body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers.
	for k, v := range u.headers {
		req.Header.Set(k, v)
	}

	// Create an HTTP client with the specified timeout.
	client := &http.Client{
		Timeout: time.Duration(u.TimeoutSec) * time.Second,
	}

	// Send the request.
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Check for non-successful status codes.
	if res.StatusCode < 200 || res.StatusCode >= 400 {
		return res, fmt.Errorf("http request failed with status code: %d %s", res.StatusCode, http.StatusText(res.StatusCode))
	}

	return res, nil
}

// authenticationHeader creates a Basic Authentication header from a "username[:password]" string.
//
// Parameters:
//
//	userPassword: The user and password string.
//
// Returns:
//
//	A tuple containing the header key ("Authorization") and the header value ("Basic <base64_token>").
func authenticationHeader(userPassword string) entity.Tuple[string, string] {
	auth := base64.StdEncoding.EncodeToString([]byte(userPassword))
	return entity.Tuple[string, string]{
		Key:   "Authorization",
		Value: "Basic " + auth,
	}
}
