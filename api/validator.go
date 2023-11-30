package api

import (
	"regexp"

	"github.com/Evans-Prah/simplebank/db/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

func validateCustomUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	// Use a regular expression to match the allowed characters
	// Adjust the regex pattern based on your requirements
	regexPattern := "^[a-zA-Z0-9_-]*$"
	matched, _ := regexp.MatchString(regexPattern, username)

	if !matched {
		return false
	}

	// Check length constraints (between 3 and 30 characters)
	return len(username) >= 3 && len(username) <= 30
}