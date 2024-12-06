package exception

type NotFoundError struct {
	message string
}

func NewNotFoundError(error string) *NotFoundError {
	return &NotFoundError{message: error}
}

func (e *NotFoundError) Error() string {
	return e.message
}
