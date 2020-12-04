// Dmitry Kisler © 2020-present
// www.dkisler.com <admin@dkisler.com>

package validator

import (
	"reflect"
	"regexp"

	validator "github.com/go-playground/validator/v10"
)

// ValidationError represents validation error content.
type ValidationError struct {
	Namespace       string
	Field           string
	StructNamespace string
	StructField     string
	Tag             string
	ActualTag       string
	Kind            reflect.Kind
	Type            reflect.Type
	Value           interface{}
	Param           string
}

// GetValidationErrors outputs validation errors.
func GetValidationErrors(err error) []ValidationError {
	var validationErrors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		var validationErr ValidationError
		validationErr.Namespace = err.Namespace()
		validationErr.Field = err.Field()
		validationErr.StructNamespace = err.StructNamespace()
		validationErr.StructField = err.StructField()
		validationErr.Tag = err.Tag()
		validationErr.ActualTag = err.ActualTag()
		validationErr.Kind = err.Kind()
		validationErr.Type = err.Type()
		validationErr.Value = err.Value()
		validationErr.Param = err.Param()
		validationErrors = append(validationErrors, validationErr)
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
