package validations

import (
	"gin-gorm-rest-api/utils"
	"mime/multipart"
)

type ValidateImageOptions struct {
	Validations []struct {
		Validate func() bool
		Message  string
		Tag      string
	}
	Field           string
	Value           *multipart.FileHeader
	ErrorValidation *[]utils.ErrorResponse
}

func ValidateImage(options ValidateImageOptions) {
	for _, validate := range options.Validations {
		if !validate.Validate() {
			*options.ErrorValidation = append(*options.ErrorValidation, utils.ErrorResponse{
				Field:   options.Field,
				Message: validate.Message,
				Tag:     validate.Tag,
				Value:   options.Value,
			})
			break
		}
	}
}
