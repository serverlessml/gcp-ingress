// Copyright 2020 dkisler.com Dmitry Kisler
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND
// NONINFRINGEMENT. IN NO EVENT WILL THE LICENSOR OR OTHER CONTRIBUTORS
// BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF, OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// See the License for the specific language governing permissions and
// limitations under the License.

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
