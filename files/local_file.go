package files

import (
	"io"
	"io/fs"
	"os"
	"strings"
)

// LocalFile is a concrete file system implementation of IFile interface
type LocalFile struct {
	uri  string
	file *os.File
}

// NewLocalFile factory method
func NewLocalFile(uri string) IFile {
	return &LocalFile{uri: uri}
}

// URI returns the resource URI with schema
// Schema can be: file, gcs, http etc
func (f *LocalFile) URI() string {
	return f.uri
}

// Close release the file resource
func (f *LocalFile) Close() error {
	if f.file != nil {
		return f.file.Close()
	} else {
		return nil
	}
}

// Read implements io.Reader interface
func (f *LocalFile) Read(p []byte) (int, error) {
	// Ensure file is open
	if err := f.openFile(); err != nil {
		return 0, err
	}
	return f.file.Read(p)
}

// Write implements io.Writer interface
func (f *LocalFile) Write(p []byte) (int, error) {
	// Ensure file is open
	if err := f.openFile(); err != nil {
		return 0, err
	}
	return f.Write(p)
}

// ReadAll read resource content to a byte array in a single call
func (f *LocalFile) ReadAll() ([]byte, error) {
	// Ensure file is open
	if err := f.openFile(); err != nil {
		return nil, err
	}

	defer f.Close()
	return io.ReadAll(f.file)
}

// WriteAll write content to a resource in a single call
func (f *LocalFile) WriteAll(b []byte) (int, error) {
	if path, err := GetUriPath(f.uri); err != nil {
		return 0, err
	} else {
		return len(b), os.WriteFile(path, b, fs.ModePerm)
	}
}

// Exists test for resource existence
func (f *LocalFile) Exists() (result bool) {
	path, err := GetUriPath(f.uri)
	if err != nil {
		return false
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// Rename change the resource name using pattern.
// The pattern can be a file or keeping parts from the original file using template ({{path}}, {{file}}, {{ext}})
func (f *LocalFile) Rename(pattern string) (string, error) {
	oldPath, err := GetUriPath(f.uri)
	if err != nil {
		return "", err
	}
	newPath := pattern

	_, path, file, ext, err := ParseUri(f.uri)
	if err != nil {
		return "", err
	}

	newPath, err = GetUriPath(newPath)
	if err != nil {
		return "", err
	}

	newPath = strings.ReplaceAll(newPath, "{{path}}", path)
	newPath = strings.ReplaceAll(newPath, "{{file}}", file)
	newPath = strings.ReplaceAll(newPath, "{{ext}}", ext)

	if err = os.Rename(oldPath, newPath); err != nil {
		return "", err
	}

	// Set current URI
	newUri := strings.ReplaceAll(f.uri, "{{path}}", path)
	newUri = strings.ReplaceAll(newUri, "{{file}}", file)
	newUri = strings.ReplaceAll(newUri, "{{ext}}", ext)
	f.uri = newUri
	f.file = nil

	return newPath, nil
}

// Delete resource
func (f *LocalFile) Delete() error {
	if path, err := GetUriPath(f.uri); err != nil {
		return err
	} else {
		return os.Remove(path)
	}
}

// Copy file content to a writer
func (f *LocalFile) Copy(wc io.WriteCloser) (int64, error) {
	// Ensure file is open
	if err := f.openFile(); err != nil {
		return 0, err
	}

	defer f.Close()
	return io.Copy(wc, f.file)
}

// Ensure file is open before any read/write operation
func (f *LocalFile) openFile() error {
	if f.file != nil {
		return nil
	}

	path, err := GetUriPath(f.uri)
	if err != nil {
		return err
	}

	// If file is not already opened, open it
	f.file, err = os.Open(path)
	return err
}
