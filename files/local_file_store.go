package files

import (
	"net/url"
	"os"
	"regexp"
	"strings"
)

// LocalFileStore is a concrete implementation of FileStore interface
type LocalFileStore struct {
	uri string
}

// NewLocalFileStore factory method
func NewLocalFileStore(uri string) IFileStore {
	return &LocalFileStore{uri: uri}
}

// URI returns the resource URI with schema
// Schema can be: file, gcs, http etc
func (f *LocalFileStore) URI() string {
	return f.uri
}

// List files in the file store
func (f *LocalFileStore) List(filter string) ([]IFile, error) {

	result := make([]IFile, 0)
	cb := func(filePath string) {
		result = append(result, NewLocalFile(filePath))
	}

	err := f.Apply(filter, cb)
	return result, err
}

// Apply action on files in the file store
func (f *LocalFileStore) Apply(filter string, action func(string)) error {

	dirPath := f.uri
	if uri, err := url.Parse(dirPath); err == nil {
		dirPath = uri.Path
	}

	rgx, er := regexp.Compile(filter)
	if er != nil {
		if len(filter) > 0 {
			return er
		}
	}

	list, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, v := range list {
		if v.IsDir() {
			continue
		}
		if fi, er := v.Info(); er == nil {
			filePath := CombineUri(f.uri, fi.Name())
			if rgx == nil {
				action(filePath)
			} else {
				if rgx.MatchString(filePath) {
					action(filePath)
				}
			}
		}
	}
	return nil
}

// Exists test for resource existence
func (f *LocalFileStore) Exists(uri string) (result bool) {
	if strings.HasPrefix(uri, "file://") {
		return NewLocalFile(uri).Exists()
	} else {
		return NewLocalFile(CombineUri(f.uri, uri)).Exists()
	}
}

// Delete resource
func (f *LocalFileStore) Delete(uri string) (err error) {
	if strings.HasPrefix(uri, "file://") {
		return NewLocalFile(uri).Delete()
	} else {
		return NewLocalFile(CombineUri(f.uri, uri)).Delete()
	}
}
