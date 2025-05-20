package entity

type Error interface {
	Error() string
	Code() int
}

// errorStruct model implements Error interface
type errorStruct struct {
	code int    // Error Code
	text string // Error Text
}

func (e *errorStruct) Error() string {
	return e.text
}

func (e *errorStruct) Code() int {
	return e.code
}

func NewError(code int, text string) Error {
	return &errorStruct{code, text}
}
