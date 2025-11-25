package files

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// HttpFile provides a concrete implementation of the IFile interface for files accessible via HTTP.
// It allows for reading from and writing to URLs, as well as checking for existence and deleting files.
type HttpFile struct {
	uri string
}

// NewHttpFile is a factory method that creates a new HttpFile instance.
//
// Parameters:
//
//	uri: The HTTP or HTTPS URL of the file.
//
// Returns:
//
//	An IFile instance representing the file at the given URI.
func NewHttpFile(uri string) IFile {
	return &HttpFile{uri: uri}
}

// URI returns the resource's URI (in this case, the URL).
//
// Returns:
//
//	The URI as a string.
func (f *HttpFile) URI() string {
	return f.uri
}

// Close releases any resources associated with the file. For HttpFile, this is a no-op
// as HTTP connections are typically managed by the net/http package on a per-request basis.
//
// Returns:
//
//	Always returns nil.
func (f *HttpFile) Close() error {
	return nil
}

// Read is part of the io.Reader interface. For HttpFile, this method is not supported
// as it's designed for streaming reads, which are better handled by ReadAll or Copy.
//
// Returns:
//
//	An error indicating that the operation is not supported.
func (f *HttpFile) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("read not supported, use ReadAll or Copy instead")
}

// Write is part of the io.Writer interface. For HttpFile, this method is not supported
// as it's designed for streaming writes, which are better handled by WriteAll.
//
// Returns:
//
//	An error indicating that the operation is not supported.
func (f *HttpFile) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("write not supported, use WriteAll instead")
}

// ReadAll reads the entire content of the file from the HTTP URL into a byte slice.
//
// Returns:
//
//	A byte slice containing the file's content.
//	An error if the HTTP GET request fails or the status code is not 200 OK.
func (f *HttpFile) ReadAll() ([]byte, error) {
	resp, err := http.Get(f.uri)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// WriteAll writes a byte slice to the file at the HTTP URL using a POST request.
//
// Parameters:
//
//	b: The byte slice to write.
//
// Returns:
//
//	The number of bytes written.
//	An error if the HTTP POST request fails.
func (f *HttpFile) WriteAll(b []byte) (int, error) {
	client := &http.Client{}
	r := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPost, f.uri, r)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// It's good practice to read the body to allow connection reuse, even if you don't use the content.
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

// Exists checks if the file exists at the HTTP URL by sending a HEAD request.
//
// Returns:
//
//	True if the server responds with a 200 OK status, false otherwise.
func (f *HttpFile) Exists() bool {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodHead, f.uri, nil)
	if err != nil {
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// Rename is not supported for HttpFile, as renaming files via HTTP is not a standard operation.
//
// Returns:
//
//	An error indicating that the operation is not supported.
func (f *HttpFile) Rename(pattern string) (string, error) {
	return "", fmt.Errorf("rename not supported for http files")
}

// Delete removes the file at the HTTP URL by sending a DELETE request.
//
// Returns:
//
//	An error if the HTTP DELETE request fails.
func (f *HttpFile) Delete() error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, f.uri, nil)
	if err != nil {
		return err
	}

	_, err = client.Do(req)
	return err
}

// Copy copies the content of the file from the HTTP URL to an io.WriteCloser.
//
// Parameters:
//
//	wc: The destination writer.
//
// Returns:
//
//	The number of bytes copied.
//	An error if the HTTP GET request fails or the copy operation fails.
func (f *HttpFile) Copy(wc io.WriteCloser) (int64, error) {
	defer wc.Close()

	resp, err := http.Get(f.uri)
	if err != nil {
		return 0, fmt.Errorf("HTTP GET error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	written, err := io.Copy(wc, resp.Body)
	return written, err
}
