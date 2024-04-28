package utils

import (
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Tag     string      `json:"tag"`
	Value   interface{} `json:"value"`
}

func Validate(dto interface{}) []ErrorResponse {
	var validate = validator.New(validator.WithRequiredStructEnabled())
	errors := []ErrorResponse{}
	err := validate.Struct(dto)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorResponse := ErrorResponse{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   err.Value(),
				Message: err.Error(),
			}
			errors = append(errors, errorResponse)
		}
	}
	return errors
}
