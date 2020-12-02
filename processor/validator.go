// Dmitry Kisler © 2020-present
// www.dkisler.com <admin@dkisler.com>

package processor

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"
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
