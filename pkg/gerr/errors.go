package gerr

type AppError interface {
	Error() string
	StatusCode() int
	UserMessage() string
	PublicCode() ErrorCode
}

type appError struct {
	message    string
	userMsg    string
	statusCode int
	publicCode ErrorCode
}

func newAppError(message, userMsg string, statusCode int, publicCode ErrorCode) AppError {
	return &appError{
		message:    message,
		userMsg:    userMsg,
		statusCode: statusCode,
		publicCode: publicCode,
	}
}

func (e *appError) Error() string {
	return e.message
}

func (e *appError) StatusCode() int {
	return e.statusCode
}

func (e *appError) UserMessage() string {
	return e.userMsg
}

func (e *appError) PublicCode() ErrorCode {
	return e.publicCode
}

func NewBusinessLogicErr(msg string, code ErrorCode) AppError {
	return newAppError(msg, msg, 422, code)
}

func NewValidationError(message string, invalidFields map[string]string) AppError {
	return &ValidationError{
		appError: appError{
			message:    message,
			userMsg:    message,
			statusCode: 400,
		},
		invalidFields: invalidFields,
	}
}

type ValidationError struct {
	appError
	invalidFields map[string]string
}

func (e *ValidationError) InvalidFields() map[string]string {
	return e.invalidFields
}
