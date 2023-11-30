package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ApiResponse struct {
	Code    int	`json:"code"`
	Message string	`json:"message"`
	Data    interface{}	`json:"data"`
	Errors    interface{}	`json:"errors,omitempty"`
}


func ApiResponseFunc(code int, message string, data interface{}, errors ...interface{}) ApiResponse {
	apiResponse := ApiResponse{
		Code:    code,
		Message: message,
		Data:    data,
		Errors: errors,
	}

	if errors != nil {
        apiResponse.Errors = errors
  }

  return apiResponse
}

func formatValidationErrors(err validator.ValidationErrors) map[string]string {
	validationErrors := make(map[string]string)

	for _, fieldError := range err {
		validationErrors[fieldError.Field()] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fieldError.Field(), fieldError.Tag())
	}

	return validationErrors
}
