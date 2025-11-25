package files

import "io"

// IFile defines the interface for a file, abstracting the underlying storage mechanism.
// This allows for concrete implementations for various file types, such as local files,
// HTTP-accessible files, or files in cloud storage like Google Cloud Storage (GCS) or AWS S3.
type IFile interface {
	// ReadWriteCloser embeds the io.ReadWriteCloser interface, providing the Read, Write, and Close methods.
	io.ReadWriteCloser

	// URI returns the Uniform Resource Identifier (URI) of the file, including the schema.
	// The schema indicates the type of file storage, e.g., "file" for local files,
	// "gcs" for Google Cloud Storage, or "http" for files accessible via HTTP.
	//
	// Returns:
	//   A string representing the URI of the file.
	URI() string

	// Exists checks if the file exists at the specified URI.
	//
	// Returns:
	//   A boolean value, true if the file exists, false otherwise.
	Exists() (result bool)

	// Rename changes the name of the file based on a given pattern.
	// The pattern can be an absolute new name or a template that uses parts of the original file name,
	// such as its path, name, or extension. For example, a pattern could be "{path}/{name}_new.{ext}".
	//
	// Parameters:
	//   pattern: The pattern to use for renaming the file.
	//
	// Returns:
	//   The new name of the file as a string.
	//   An error if the renaming process fails.
	Rename(pattern string) (result string, err error)

	// Delete removes the file from its storage.
	//
	// Returns:
	//   An error if the deletion fails.
	Delete() (err error)

	// ReadAll reads the entire content of the file into a byte slice.
	// This is a convenience method for reading the whole file in a single operation.
	//
	// Returns:
	//   A byte slice containing the file's content.
	//   An error if reading fails.
	ReadAll() (b []byte, err error)

	// WriteAll writes a byte slice to the file, overwriting any existing content.
	// This is a convenience method for writing the entire content in a single operation.
	//
	// Parameters:
	//   b: The byte slice to write to the file.
	//
	// Returns:
	//   The number of bytes written.
	//   An error if writing fails.
	WriteAll(b []byte) (n int, err error)

	// Copy copies the content of the file to an io.WriteCloser.
	// This is useful for streaming the file's content to another destination.
	//
	// Parameters:
	//   wc: The io.WriteCloser to which the file content will be copied.
	//
	// Returns:
	//   The total number of bytes written.
	//   An error if the copy operation fails.
	Copy(wc io.WriteCloser) (written int64, err error)
}
