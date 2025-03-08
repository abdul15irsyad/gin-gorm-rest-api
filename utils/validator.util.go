package utils

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   any    `json:"value"`
}

func Validate[T any](dtos T) []ErrorResponse {
	var validate = validator.New()
	errors := []ErrorResponse{}
	err := validate.Struct(dtos)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.StructField()
			field, _ := reflect.TypeOf(dtos).FieldByName(fieldName)
			tagValue := field.Tag.Get("json")
			if tagValue == "" {
				tagValue = field.Tag.Get("form")
			}
			if tagValue == "" {
				tagValue = err.Field()
			}
			errorResponse := ErrorResponse{
				Field:   tagValue,
				Tag:     err.Tag(),
				Value:   err.Value(),
				Message: err.Error(),
			}
			errors = append(errors, errorResponse)
		}
	}
	return errors
}
