package files

import (
	"fmt"
	"net/url"
	"strings"
)

// IFileStore defines the interface for a file store, which is a repository for files.
// This interface provides methods for listing, managing, and interacting with files
// in a way that is agnostic to the underlying storage system (e.g., local file system,
// HTTP server, Google Cloud Storage, AWS S3).
type IFileStore interface {
	// URI returns the base Uniform Resource Identifier (URI) of the file store,
	// including the schema (e.g., "file", "gcs", "http").
	//
	// Returns:
	//   A string representing the URI of the file store.
	URI() string

	// List retrieves a list of files within the store that match a given filter.
	// The filter is typically a regular expression used to match file names.
	//
	// Parameters:
	//   filter: A regular expression string to filter the files.
	//
	// Returns:
	//   A slice of IFile instances that match the filter.
	//   An error if the listing operation fails.
	List(filter string) (result []IFile, err error)

	// Apply executes a given action on each file in the file store that matches a filter.
	// The action is a function that takes the file's URI as a string.
	//
	// Parameters:
	//   filter: A regular expression string to filter the files.
	//   action: A function to apply to each matching file's URI.
	//
	// Returns:
	//   An error if the apply operation fails.
	Apply(filter string, action func(string)) error

	// Exists checks if a file with the specified URI exists within the store.
	//
	// Parameters:
	//   uri: The URI of the file to check.
	//
	// Returns:
	//   A boolean value, true if the file exists, false otherwise.
	Exists(uri string) (result bool)

	// Delete removes a file with the specified URI from the store.
	//
	// Parameters:
	//   uri: The URI of the file to delete.
	//
	// Returns:
	//   An error if the deletion fails.
	Delete(uri string) (err error)

	// Close releases any resources associated with the file store, such as network connections.
	//
	// Returns:
	//   An error if closing fails.
	Close() error
}

// CombineUri constructs a URI by joining a base URI with one or more path segments.
// It ensures that there are no trailing slashes in the base URI or parts before joining.
//
// Parameters:
//
//	uri: The base URI.
//	parts: A variadic slice of strings representing the path segments to append.
//
// Returns:
//
//	A single string representing the combined URI.
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

// GetUriPath extracts the path component from a given URI.
//
// Parameters:
//
//	uri: The URI string to parse.
//
// Returns:
//
//	The path component of the URI as a string.
//	An error if the URI cannot be parsed.
func GetUriPath(uri string) (path string, err error) {
	Url, er := url.Parse(uri)
	if er != nil {
		return "", er
	}
	return Url.Path, nil
}

// ParseUri decomposes a URI into its constituent parts: schema, path, file name, and extension.
//
// Parameters:
//
//	uri: The URI string to parse.
//
// Returns:
//
//	The schema (e.g., "http", "file").
//	The path of the resource.
//	The file name without the extension.
//	The file extension.
//	An error if the URI cannot be parsed.
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
