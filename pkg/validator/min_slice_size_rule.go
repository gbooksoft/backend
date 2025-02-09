package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func minSliceSizeValidator(fl validator.FieldLevel) bool {
	field := fl.Field()

	if field.Kind() != reflect.Slice {
		return false
	}

	minSize := fl.Param()

	var minLength int
	fmt.Sscanf(minSize, "%d", &minLength)

	return field.Len() >= minLength
}
