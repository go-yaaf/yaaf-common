package files

import (
	"fmt"
	"net/url"
	"strings"
)

// IFileStore Files store interface
// This interface is used for concrete implementation of any file store (local file system, HTTP files, Google / AWS Buckets etc.
type IFileStore interface {

	// URI returns the resource URI with schema
	// Schema can be: file, gcs, http etc
	URI() string

	// List files in the URI using regexp filter
	List(filter string) (result []IFile, err error)

	// Apply action on files in the file store
	Apply(filter string, action func(string)) error

	// Exists test for resource existence
	Exists(uri string) (result bool)

	// Delete resource
	Delete(uri string) (err error)

	// Release assosiated resources, if any
	Close() error
}

// CombineUri creates a URI from segments
func CombineUri(uri string, parts ...string) string {
	if strings.HasSuffix(uri, "/") {
		uri = uri[:len(uri)-1]
	}
	result := uri

	for _, part := range parts {
		if strings.HasSuffix(part, "/") {
			part = part[:len(part)-1]
		}
		result = fmt.Sprintf("%s/%s", result, part)
	}
	return result
}

// GetUriPath extract path from the uri
func GetUriPath(uri string) (path string, err error) {
	if Url, er := url.Parse(uri); er != nil {
		return "", er
	} else {
		return Url.Path, nil
	}
}

// ParseUri decompose URI to parts
func ParseUri(uri string) (schema, path, file, ext string, err error) {

	Url, er := url.Parse(uri)
	if er != nil {
		return "", "", "", "", er
	}

	schema = Url.Scheme
	idx := strings.LastIndex(Url.Path, "/")
	if idx < 0 {
		path = "/"
	} else {
		path = Url.Path[0:idx]
	}

	ide := strings.LastIndex(Url.Path, ".")
	if ide == -1 {
		file = Url.Path[idx+1:]
	} else {
		file = Url.Path[idx+1 : ide]
		ext = Url.Path[ide+1:]
	}
	return
}
