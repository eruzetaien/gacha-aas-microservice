package exception

type UserError struct {
	message string
}

func NewUserError(error string) *UserError {
	return &UserError{message: error}
}

func (e *UserError) Error() string {
	return e.message
}
