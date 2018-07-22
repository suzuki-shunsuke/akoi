package util

import (
	"gopkg.in/go-playground/validator.v9"
)

var (
	validate = validator.New()
)

// ValidateStruct validates a struct.
func ValidateStruct(v interface{}) error {
	return validate.Struct(v)
}
