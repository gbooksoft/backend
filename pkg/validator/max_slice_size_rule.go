package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func maxSliceSizeValidator(fl validator.FieldLevel) bool {
	field := fl.Field()

	if field.Kind() != reflect.Slice {
		return false
	}

	maxSize := fl.Param()

	var maxLength int
	fmt.Sscanf(maxSize, "%d", &maxLength)

	return field.Len() <= maxLength
}
