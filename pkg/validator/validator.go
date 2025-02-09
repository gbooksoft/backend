package validator

import "github.com/go-playground/validator/v10"

func NewCustomValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("min_slice_size", minSliceSizeValidator)
	validate.RegisterValidation("max_slice_size", maxSliceSizeValidator)
	validate.RegisterValidation("uuid4_slice", uuid4SliceValidator)

	return validate
}
