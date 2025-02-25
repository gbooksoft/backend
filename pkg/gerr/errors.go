package gerr

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

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

func NewValidationError(err error) AppError {
	messages := make(map[string]string)

	errors, ok := err.(validator.ValidationErrors)

	if ok {
		for _, err := range errors {
			messages[err.Field()] = err.Error()
		}
	} else {
		log.Error().Err(err).Msg("Error parsing validation error")
		messages["unknown"] = err.Error()
	}

	return &ValidationError{
		appError: appError{
			message:    "invalid input",
			userMsg:    "invalid input",
			statusCode: 400,
		},
		invalidFields: messages,
	}
}

type ValidationError struct {
	appError
	invalidFields map[string]string
}

func (e *ValidationError) InvalidFields() map[string]string {
	return e.invalidFields
}

func NewNotFoundErr(msg string) AppError {
	return newAppError(msg, msg, 404, NotFound)
}
