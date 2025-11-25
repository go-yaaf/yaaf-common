package files

import "fmt"

// HttpFileStore provides a concrete implementation of the IFileStore interface for files accessible via HTTP.
// While it implements the interface, some operations like listing and applying actions are not supported
// due to the nature of HTTP, which doesn't typically provide a standard way to list directory-like contents.
type HttpFileStore struct {
	uri string
}

// NewHttpFileStore is a factory method that creates a new HttpFileStore instance.
//
// Parameters:
//
//	uri: The base HTTP or HTTPS URL for the file store.
//
// Returns:
//
//	An IFileStore instance.
func NewHttpFileStore(uri string) IFileStore {
	return &HttpFileStore{uri: uri}
}

// URI returns the base URI of the file store.
//
// Returns:
//
//	The base URI as a string.
func (f *HttpFileStore) URI() string {
	return f.uri
}

// List is not supported for HttpFileStore. HTTP does not have a standard protocol
// for listing files in a directory-like manner.
//
// Returns:
//
//	An error indicating that the operation is not supported.
func (f *HttpFileStore) List(filter string) ([]IFile, error) {
	return nil, fmt.Errorf("list operation is not supported for http file stores")
}

// Apply is not supported for HttpFileStore. Since listing files is not supported,
// applying an action to a set of files is also not feasible.
//
// Returns:
//
//	An error indicating that the operation is not supported.
func (f *HttpFileStore) Apply(filter string, action func(string)) error {
	return fmt.Errorf("apply operation is not supported for http file stores")
}

// Exists checks if a file at a given URI exists. It delegates this check to an HttpFile instance.
//
// Parameters:
//
//	uri: The full URI of the file to check.
//
// Returns:
//
//	True if the file exists, false otherwise.
func (f *HttpFileStore) Exists(uri string) bool {
	return NewHttpFile(uri).Exists()
}

// Delete removes a file at a given URI. It delegates the deletion to an HttpFile instance.
//
// Parameters:
//
//	uri: The full URI of the file to delete.
//
// Returns:
//
//	An error if the deletion fails.
func (f *HttpFileStore) Delete(uri string) error {
	return NewHttpFile(uri).Delete()
}

// Close releases any resources held by the file store. For HttpFileStore, this is a no-op.
//
// Returns:
//
//	Always returns nil.
func (f *HttpFileStore) Close() error {
	return nil
}
