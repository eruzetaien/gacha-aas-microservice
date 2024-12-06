package exception

type ConflictError struct {
	message string
}

func NewConflictError(error string) *ConflictError {
	return &ConflictError{message: error}
}

func (e *ConflictError) Error() string {
	return e.message
}
