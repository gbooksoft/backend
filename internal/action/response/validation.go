package response

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

func parseValidationErr(err error) map[string]string {
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

	return messages
}
