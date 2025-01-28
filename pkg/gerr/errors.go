package gerr

type BusinessLogicErr struct {
	Msg  string
	Code ErrorCode
	err  error
}

func NewBusinessLogicErr(msg string, code ErrorCode, err error) BusinessLogicErr {
	return BusinessLogicErr{Msg: msg, Code: code, err: err}
}

func (err BusinessLogicErr) Error() string {
	msg := err.Msg
	if err.err != nil {
		msg = err.err.Error()
	}

	return "business logic error: " + msg
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
