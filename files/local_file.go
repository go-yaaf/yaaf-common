package files

import (
	"io"
	"io/fs"
	"os"
	"strings"
)

// LocalFile provides a concrete implementation of the IFile interface for files on the local file system.
// It wraps the standard `os.File` and provides methods for interacting with files using URIs.
type LocalFile struct {
	uri  string
	file *os.File
}

// NewLocalFile is a factory method that creates a new LocalFile instance.
//
// Parameters:
//
//	uri: The URI of the local file, which should have the "file" schema (e.g., "file:///path/to/file").
//
// Returns:
//
//	An IFile instance representing the local file.
func NewLocalFile(uri string) IFile {
	return &LocalFile{uri: uri}
}

// URI returns the resource's URI.
//
// Returns:
//
//	The URI as a string.
func (f *LocalFile) URI() string {
	return f.uri
}

// Close releases the file resource by closing the underlying `os.File`.
// If the file is not open, it does nothing.
//
// Returns:
//
//	An error if closing the file fails.
func (f *LocalFile) Close() error {
	if f.file != nil {
		err := f.file.Close()
		f.file = nil // Make sure to nullify the file handle after closing.
		return err
	}
	return nil
}

// Read reads up to len(p) bytes into p. It returns the number of bytes read and any error encountered.
// It ensures the file is open before reading.
//
// Returns:
//
//	The number of bytes read.
//	An error if the file cannot be opened or if the read operation fails.
func (f *LocalFile) Read(p []byte) (int, error) {
	if err := f.openFile(os.O_RDONLY, 0); err != nil {
		return 0, err
	}
	return f.file.Read(p)
}

// Write writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p)) and any error encountered that caused the write to stop early.
// It ensures the file is open for writing before the operation.
//
// Returns:
//
//	The number of bytes written.
//	An error if the file cannot be opened or if the write operation fails.
func (f *LocalFile) Write(p []byte) (int, error) {
	if err := f.openFile(os.O_WRONLY|os.O_CREATE, 0666); err != nil {
		return 0, err
	}
	return f.file.Write(p)
}

// ReadAll reads the entire content of the file into a byte slice.
//
// Returns:
//
//	A byte slice containing the file's content.
//	An error if reading fails.
func (f *LocalFile) ReadAll() ([]byte, error) {
	path, err := GetUriPath(f.uri)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

// WriteAll writes a byte slice to the file, creating it if necessary and overwriting existing content.
//
// Parameters:
//
//	b: The byte slice to write.
//
// Returns:
//
//	The number of bytes written.
//	An error if writing fails.
func (f *LocalFile) WriteAll(b []byte) (int, error) {
	path, err := GetUriPath(f.uri)
	if err != nil {
		return 0, err
	}
	err = os.WriteFile(path, b, fs.ModePerm)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

// Exists checks if the file exists on the local file system.
//
// Returns:
//
//	True if the file exists, false otherwise.
func (f *LocalFile) Exists() bool {
	path, err := GetUriPath(f.uri)
	if err != nil {
		return false
	}

	_, err = os.Stat(path)
	return !os.IsNotExist(err)
}

// Rename changes the name of the file using a specified pattern.
// The pattern can include placeholders like {{path}}, {{file}}, and {{ext}} to reuse parts of the original name.
//
// Parameters:
//
//	pattern: The new name pattern.
//
// Returns:
//
//	The new path of the file.
//	An error if renaming fails.
func (f *LocalFile) Rename(pattern string) (string, error) {
	oldPath, err := GetUriPath(f.uri)
	if err != nil {
		return "", err
	}

	_, path, file, ext, err := ParseUri(f.uri)
	if err != nil {
		return "", err
	}

	newPath := strings.ReplaceAll(pattern, "{{path}}", path)
	newPath = strings.ReplaceAll(newPath, "{{file}}", file)
	newPath = strings.ReplaceAll(newPath, "{{ext}}", ext)

	if err = os.Rename(oldPath, newPath); err != nil {
		return "", err
	}

	// Update the URI to reflect the new path.
	f.uri = "file://" + newPath
	f.file = nil // The old file handle is no longer valid.
	return newPath, nil
}

// Delete removes the file from the local file system.
//
// Returns:
//
//	An error if the deletion fails.
func (f *LocalFile) Delete() error {
	path, err := GetUriPath(f.uri)
	if err != nil {
		return err
	}
	return os.Remove(path)
}

// Copy copies the content of the file to an io.WriteCloser.
//
// Parameters:
//
//	wc: The destination writer.
//
// Returns:
//
//	The number of bytes copied.
//	An error if the copy operation fails.
func (f *LocalFile) Copy(wc io.WriteCloser) (int64, error) {
	defer wc.Close()
	if err := f.openFile(os.O_RDONLY, 0); err != nil {
		return 0, err
	}
	defer f.Close()
	return io.Copy(wc, f.file)
}

// openFile ensures that the file is open before any read or write operation.
// It opens the file with the specified flag and permission.
func (f *LocalFile) openFile(flag int, perm fs.FileMode) error {
	if f.file != nil {
		return nil
	}

	path, err := GetUriPath(f.uri)
	if err != nil {
		return err
	}

	f.file, err = os.OpenFile(path, flag, perm)
	return err
}
