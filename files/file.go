package files

import "io"

// IFile File interface
// This interface is used for concrete implementation of any file (local file system, HTTP files, Google / AWS Buckets etc).
type IFile interface {

	// ReadWriteCloser must implement Read(), Write() and Close() method
	io.ReadWriteCloser

	// URI returns the resource URI with schema
	// Schema can be: file, gcs, http etc
	URI() string

	// Exists test for resource existence
	Exists() (result bool)

	// Rename change the resource name using pattern.
	// The pattern can be an absolute name or keeping parts from the original file using template (path, name, ext)
	Rename(pattern string) (result string, err error)

	// Delete resource
	Delete() (err error)

	// ReadAll read resource content to a byte array in a single call
	ReadAll() (b []byte, err error)

	// WriteAll write content to a resource in a single call
	WriteAll(b []byte) (n int, err error)

	// Copy file content to a writer
	Copy(wc io.WriteCloser) (written int64, err error)
}
