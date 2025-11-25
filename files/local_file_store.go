package files

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// LocalFileStore provides a concrete implementation of the IFileStore interface for the local file system.
// It allows for listing and managing files within a specified directory.
type LocalFileStore struct {
	uri string
}

// NewLocalFileStore is a factory method that creates a new LocalFileStore instance.
//
// Parameters:
//
//	uri: The URI of the base directory for the file store, which should have the "file" schema (e.g., "file:///path/to/directory").
//
// Returns:
//
//	An IFileStore instance representing the local file store.
func NewLocalFileStore(uri string) IFileStore {
	return &LocalFileStore{uri: uri}
}

// URI returns the base URI of the file store.
//
// Returns:
//
//	The base URI as a string.
func (f *LocalFileStore) URI() string {
	return f.uri
}

// List retrieves a list of files within the store that match a given regex filter.
//
// Parameters:
//
//	filter: A regular expression to match against file paths.
//
// Returns:
//
//	A slice of IFile instances that match the filter.
//	An error if the listing operation fails.
func (f *LocalFileStore) List(filter string) ([]IFile, error) {
	var files []IFile
	action := func(filePath string) {
		files = append(files, NewLocalFile(filePath))
	}

	if err := f.Apply(filter, action); err != nil {
		return nil, err
	}
	return files, nil
}

// Apply executes a given action on each file in the file store that matches a regex filter.
//
// Parameters:
//
//	filter: A regular expression to match against file paths.
//	action: A function to apply to each matching file's URI.
//
// Returns:
//
//	An error if the apply operation fails.
func (f *LocalFileStore) Apply(filter string, action func(string)) error {
	dirPath, err := f.getDirPath()
	if err != nil {
		return err
	}

	regex, err := regexp.Compile(filter)
	if err != nil {
		return fmt.Errorf("invalid filter regex: %w", err)
	}

	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && regex.MatchString(path) {
			action("file://" + path)
		}
		return nil
	})
}

// Exists checks if a file exists within the store. The URI can be absolute or relative to the store's base URI.
//
// Parameters:
//
//	uri: The URI of the file to check.
//
// Returns:
//
//	True if the file exists, false otherwise.
func (f *LocalFileStore) Exists(uri string) bool {
	fullUri := f.resolveUri(uri)
	return NewLocalFile(fullUri).Exists()
}

// Delete removes a file from the store. The URI can be absolute or relative to the store's base URI.
//
// Parameters:
//
//	uri: The URI of the file to delete.
//
// Returns:
//
//	An error if the deletion fails.
func (f *LocalFileStore) Delete(uri string) error {
	fullUri := f.resolveUri(uri)
	return NewLocalFile(fullUri).Delete()
}

// Close releases any resources held by the file store. For LocalFileStore, this is a no-op.
//
// Returns:
//
//	Always returns nil.
func (f *LocalFileStore) Close() error {
	return nil
}

// getDirPath extracts the directory path from the store's URI.
func (f *LocalFileStore) getDirPath() (string, error) {
	parsedUrl, err := url.Parse(f.uri)
	if err != nil {
		return "", fmt.Errorf("invalid store URI: %w", err)
	}
	return parsedUrl.Path, nil
}

// resolveUri resolves a given URI against the store's base URI.
// If the given URI is absolute, it's returned as is. Otherwise, it's combined with the store's base URI.
func (f *LocalFileStore) resolveUri(uri string) string {
	if strings.HasPrefix(uri, "file://") {
		return uri
	}
	return CombineUri(f.uri, uri)
}
