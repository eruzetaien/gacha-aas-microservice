package exception

type BadRequestError struct {
	message string
}

func NewBadRequestError(error string) *BadRequestError {
	return &BadRequestError{message: error}
}

func (e *BadRequestError) Error() string {
	return e.message
}
