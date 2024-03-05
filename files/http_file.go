package files

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// HttpFile is a concrete http implementation of IFile interface
type HttpFile struct {
	uri string
}

// NewHttpFile factory method
func NewHttpFile(uri string) IFile {
	return &HttpFile{uri: uri}
}

// URI returns the resource URI with schema
// Schema can be: file, gcs, http etc
func (f *HttpFile) URI() string {
	return f.uri
}

// Close release the file resource
func (f *HttpFile) Close() error {
	return nil
}

// Read implements io.Reader interface
func (f *HttpFile) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("not supported")
}

// Write implements io.Writer interface
func (f *HttpFile) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("not supported")
}

// ReadAll read resource content to a byte array in a single call
func (f *HttpFile) ReadAll() ([]byte, error) {
	resp, err := http.Get(f.uri)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status error: %v", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// WriteAll write content to a resource in a single call
func (f *HttpFile) WriteAll(b []byte) (int, error) {
	client := &http.Client{}

	r := bytes.NewReader(b)
	req, err := http.NewRequest("POST", f.uri, r)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	if _, err = io.ReadAll(resp.Body); err != nil {
		return 0, err
	} else {
		return len(b), nil
	}
}

// Exists test for resource existence
func (f *HttpFile) Exists() (result bool) {
	client := &http.Client{}

	req, err := http.NewRequest("HEAD", f.uri, nil)
	if err != nil {
		return false
	}
	if res, err := client.Do(req); err != nil {
		return false
	} else {
		return res.StatusCode == 200
	}
}

// Rename change the resource name using pattern.
// The pattern can be a file or keeping parts from the original file using template ({{path}}, {{file}}, {{ext}})
func (f *HttpFile) Rename(pattern string) (string, error) {
	return pattern, fmt.Errorf("not supported")
}

// Delete resource
func (f *HttpFile) Delete() error {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", f.uri, nil)
	if err != nil {
		return err
	}
	if _, err := client.Do(req); err != nil {
		return err
	} else {
		return nil
	}
}

// Copy file content to a writer
func (f *HttpFile) Copy(wc io.WriteCloser) (int64, error) {
	resp, err := http.Get(f.uri)
	if err != nil {
		return 0, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status error: %v", resp.StatusCode)
	}

	written, err := io.Copy(wc, resp.Body)
	_ = wc.Close()
	return written, err
}
