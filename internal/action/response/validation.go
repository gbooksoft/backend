package response

import (
	"github.com/go-playground/validator/v10"
)

func parseValidationErr(err error) map[string]string {
	messages := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		messages[err.Field()] = err.Error()
	}

	return messages
}
