/*
Copyright © 2020 Dmitry Kisler <admin@dkisler.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handlers_test

import (
	"testing"

	"github.com/serverlessml/ingress/handlers"
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
