// JWT token utilities
//

package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/go-yaaf/yaaf-common/entity"
	"net/http"
	"net/url"
)

type HttpUtilsStruct struct {
	method     string
	url        string
	body       string
	headers    map[string]string
	TimeoutSec int
}

// HttpUtils is a factory method that acts as a static member
func HttpUtils() *HttpUtilsStruct {
	return &HttpUtilsStruct{
		method:  "GET",
		headers: make(map[string]string),
	}
}

func (u *HttpUtilsStruct) New(method, url string) *HttpUtilsStruct {
	u.method = method
	u.url = url
	return u
}

func (u *HttpUtilsStruct) WithHeader(key, value string) *HttpUtilsStruct {
	u.headers[key] = value
	return u
}

func (u *HttpUtilsStruct) WithHeaders(headers map[string]string) *HttpUtilsStruct {
	for k, v := range headers {
		u.headers[k] = v
	}
	return u
}

func (u *HttpUtilsStruct) WithBody(body string) *HttpUtilsStruct {
	u.body = body
	return u
}

func (u *HttpUtilsStruct) WithTimeout(timeout int) *HttpUtilsStruct {
	u.TimeoutSec = timeout
	return u
}

func (u *HttpUtilsStruct) Send() (*http.Response, error) {

	parsedUrl, err := url.Parse(u.url)
	if err != nil {
		return nil, err
	}

	// Add auth header
	if parsedUrl.User != nil {
		result := authenticationHeader(parsedUrl.User.String())
		u.headers[result.Key] = result.Value
	}

	// re-build URL
	realUrl := url.URL{
		Scheme:   parsedUrl.Scheme,
		Host:     parsedUrl.Host,
		Path:     parsedUrl.Path,
		RawQuery: parsedUrl.Query().Encode(),
	}

	// Send HTTP
	var (
		req *http.Request
		res *http.Response
	)

	defer func() {
		if res != nil {
			if res.Body != nil {
				_ = res.Body.Close()
			}
		}
	}()

	if req, err = http.NewRequest(u.method, realUrl.String(), bytes.NewBuffer([]byte(u.body))); err != nil {
		return nil, err
	}

	req.Close = true
	for k, v := range u.headers {
		req.Header.Set(k, v)
	}

	if res, err = http.DefaultClient.Do(req); err != nil {
		return nil, err
	}

	code := res.StatusCode
	if code < 200 || code >= 400 {
		if st := http.StatusText(res.StatusCode); len(st) == 0 {
			err = fmt.Errorf("http status code: %d", res.StatusCode)
		}
	}
	return res, err
}

// authenticationHeader receives  username and password information in the standard form
// of "username[:password]".
func authenticationHeader(userPassword string) entity.Tuple[string, string] {
	auth := base64.StdEncoding.EncodeToString([]byte(userPassword))
	return entity.Tuple[string, string]{
		Key:   "Authorization",
		Value: "Basic " + auth,
	}
}

// endregion
