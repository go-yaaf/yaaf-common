package files

import "fmt"

// HttpFileStore is a concrete implementation of FileStore interface
type HttpFileStore struct {
	uri string
}

// NewHttpFileStore factory method
func NewHttpFileStore(uri string) IFileStore {
	return &HttpFileStore{uri: uri}
}

// URI returns the resource URI with schema
func (f *HttpFileStore) URI() string {
	return f.uri
}

// List files in the file store
func (f *HttpFileStore) List(filter string) ([]IFile, error) {
	return nil, fmt.Errorf("not supported")
}

// Apply action on files in the file store
func (f *HttpFileStore) Apply(filter string, action func(string)) error {
	return fmt.Errorf("not supported")
}

// Exists test for resource existence
func (f *HttpFileStore) Exists(uri string) (result bool) {
	return NewHttpFile(uri).Exists()
}

// Delete resource
func (f *HttpFileStore) Delete(uri string) (err error) {
	return NewHttpFile(uri).Delete()
}
