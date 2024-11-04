package gerr

type BusinessLogicErr struct {
	Msg string
	err error
}

func NewBusinessLogicErr(msg string, err error) BusinessLogicErr {
	return BusinessLogicErr{Msg: msg, err: err}
}

func (err BusinessLogicErr) Error() string {
	return "business logic error: " + err.err.Error()
}

type ValidationError struct {
	err error
}

func NewValidationError(err error) *ValidationError {
	return &ValidationError{err: err}
}

func (err ValidationError) Error() string {
	return "validation error: " + err.err.Error()
}
