package validator

import "github.com/go-playground/validator/v10"

func NewCustomValidator() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())
}
