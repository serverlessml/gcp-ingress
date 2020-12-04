// Dmitry Kisler © 2020-present
// www.dkisler.com <admin@dkisler.com>

package validator

import (
	"regexp"

	validator "github.com/go-playground/validator/v10"
)

// GetValidationErrors outputs validation errors.
func GetValidationErrors(err error) []string {
	if err == nil {
		return []string{}
	}
	var validationErrors []string
	for _, err := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, err.Error())
	}
	return validationErrors
}

var sha1 = regexp.MustCompile(`^[a-fA-F0-9]{40}$`)

// IsSHA1 check if the input string is a valida SHA1 hash.
func IsSHA1(fl validator.FieldLevel) bool {
	return sha1.MatchString(fl.Field().String())
}

// New is instantiating a validator instance.
func New() *validator.Validate {
	Validator := validator.New()
	Validator.RegisterValidation("sha1", IsSHA1, true)
	return Validator
}
