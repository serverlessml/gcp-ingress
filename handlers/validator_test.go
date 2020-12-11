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

package handlers_test

import (
	"testing"

	"github.com/serverlessml/gcp-ingress/handlers"
)

func TestValidator(t *testing.T) {
	schema := `{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type": "number", "minimum": 1, "maximum": 2
	}`

	tests := []struct {
		name    string
		in      []byte
		want    string
		isError bool
	}{
		{
			name:    "Positive",
			in:      []byte(`1.1`),
			want:    "",
			isError: false,
		},
		{
			name:    "Parsing Error",
			in:      []byte(`{`),
			want:    "parsing",
			isError: true,
		},
		{
			name:    "Validation Error",
			in:      []byte(`10`),
			want:    "validation",
			isError: true,
		},
	}
	for _, test := range tests {
		got := handlers.Validate(schema, test.in)
		eType := ""
		if test.isError {
			eType = (got.(handlers.Error)).Type
		}
		if eType != test.want {
			t.Fatalf("[%s]: Wrong error implementation", test.name)
		}
	}
}
