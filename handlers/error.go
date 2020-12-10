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

package handlers

import (
	"encoding/json"
	"fmt"
)

// Error defines the error.
type Error struct {
	Type    string
	Message string
	Details interface{}
}

// Error returns a human readable error message.
func (e Error) Error() string {
	prefix := fmt.Sprintf("[%s]", e.Type)
	if e.Details == nil {
		return fmt.Sprintf("%s %s", prefix, e.Message)
	}
	return fmt.Sprintf("%s %s. Details:\n%v", prefix, e.Message, e.Details)
}

// ErrorPush defines the errors of submitting jobs.
type ErrorPush struct {
	// contains error message
	Message string `json:"message"`
	// contains pipeline config
	Details interface{} `json:"details"`
}

// Error returns a human readable error message.
func (e ErrorPush) Error() string {
	if e.Details == nil {
		return e.Message
	}
	return fmt.Sprintf("%s. Details:\n%v", e.Message, e.Details)
}

// ErrorArray returns human readable error messages.
func ErrorArray(errors []error) []string {
	e := []string{}
	for _, err := range errors {
		e = append(e, err.Error())
	}
	return e
}

// NewUnmarshallerError reads the json unmarshal error.
func NewUnmarshallerError(unmarshalErr error) error {
	err := unmarshalErr.(*json.UnmarshalTypeError)
	return Error{
		Type:    "parsing",
		Message: "Input parsing error",
		Details: fmt.Sprintf("Cannot unmarshal value: %s of type %s", err.Field, err.Type.String()),
	}
}
