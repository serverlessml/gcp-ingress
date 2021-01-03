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

func TestError(t *testing.T) {
	tests := []struct {
		name string
		in   handlers.Error
		want string
	}{
		{
			name: "With details",
			in: handlers.Error{
				Type:    "test",
				Message: "foo",
				Details: map[string]interface{}{"foo": "bar"},
			},
			want: "[test] foo. Details:\nmap[foo:bar]",
		},
		{
			name: "Without details",
			in: handlers.Error{
				Type:    "test",
				Message: "foo",
			},
			want: "[test] foo",
		},
	}

	for _, test := range tests {
		got := test.in.Error()
		if got != test.want {
			t.Fatalf("[%s]: Results don't match\nwant: %s\ngot: %s",
				test.name, test.want, got)
		}
	}
}

func TestErrorPush(t *testing.T) {
	tests := []struct {
		name string
		in   handlers.ErrorPush
		want string
	}{
		{
			name: "With details",
			in: handlers.ErrorPush{
				Message: "foo",
				Details: map[string]interface{}{"foo": "bar"},
			},
			want: "foo. Details:\nmap[foo:bar]",
		},
		{
			name: "Without details",
			in: handlers.ErrorPush{
				Message: "foo",
			},
			want: "foo",
		},
	}

	for _, test := range tests {
		got := test.in.Error()
		if got != test.want {
			t.Fatalf("[%s]: Results don't match\nwant: %s\ngot: %s",
				test.name, test.want, got)
		}
	}
}
