package entity

// Error interface defines the contract for error handling in the application.
// It extends the standard error interface with a numeric code, allowing for more structured error handling.
type Error interface {
	Error() string
	Code() int
}

// errorStruct is the default implementation of the Error interface.
// @Data
type errorStruct struct {
	code int    // code is the numeric error code
	text string // text is the error message
}

// Error returns the error text message.
func (e *errorStruct) Error() string {
	return e.text
}

// Code returns the numeric error code.
func (e *errorStruct) Code() int {
	return e.code
}

// NewError creates a new Error instance with the specified code and text.
func NewError(code int, text string) Error {
	return &errorStruct{code, text}
}
