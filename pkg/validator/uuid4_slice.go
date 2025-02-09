package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"reflect"
)

func uuid4SliceValidator(fl validator.FieldLevel) bool {
	field := fl.Field()

	if field.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < field.Len(); i++ {
		element := field.Index(i).String()
		_, err := uuid.Parse(element)
		if err != nil {
			return false
		}
	}

	return true
}
