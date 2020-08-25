package web

// ErrorResponse contains the message that will be sent to the client in case of error
type ErrorResponse struct{
	Error string `json:"error"`
}

// Error contains the http status code along with the error message
type Error struct {
	Err error
	StatusCode int
}

// NewRequestError is a factory method that creates Error struct
func NewRequestError(err error, statusCode int) error {
	return &Error{err, statusCode}
}

func (e *Error) Error() string {
	return e.Err.Error()
}