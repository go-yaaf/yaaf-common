package entity

// Error model extends the standard error interface with code
type Error struct {
	Code int    `json:"code"` // Error Code
	Text string `json:"text"` // Error Text
}

func (e Error) Error() string {
	return e.Text
}
